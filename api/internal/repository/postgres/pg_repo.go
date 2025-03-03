package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"lms/internal/entity"
	"lms/internal/repository"

	_ "github.com/lib/pq"
)

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *postgresRepo {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) GetAllGroups(ctx context.Context) ([]entity.Group, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM groups order by sort`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]entity.Group, 0, 10)
	for rows.Next() {
		var g entity.Group
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

func (r *postgresRepo) GetCriteriaByGroup(ctx context.Context, groupId int) ([]entity.Criteria, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT c.id, c.name, c.max_score FROM criteria c inner join group_criteria gc on gc.criteria_id = c.id where gc.group_id = $1`, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]entity.Criteria, 0, 10)
	for rows.Next() {
		var e entity.Criteria
		if err := rows.Scan(&e.ID, &e.Name, &e.MaxScore); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

func (r *postgresRepo) GetStudentsByGroup(ctx context.Context, groupId int) ([]entity.Student, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM students where group_id = $1`, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]entity.Student, 0, 10)
	for rows.Next() {
		var e entity.Student
		if err := rows.Scan(&e.ID, &e.Name); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

func (r *postgresRepo) GetLessonsWithScores(ctx context.Context, groupId int) ([]repository.LessonWithScores, error) {
	var lessonIds []int
	lessons := make([]repository.LessonWithScores, 0)

	students, err := r.GetStudentsByGroup(ctx, groupId)

	rowsLessons, err := r.db.QueryContext(ctx, `
        SELECT id, dates
          FROM lessons
         WHERE group_id = $1
        ORDER BY id DESC
    `, groupId)
	if err != nil {
		return nil, err
	}
	defer rowsLessons.Close()

	for rowsLessons.Next() {
		var lesson entity.Lesson
		var datesBytes []byte
		if err := rowsLessons.Scan(&lesson.ID, &datesBytes); err != nil {
			return nil, err
		}

		dates := make([]string, 0)
		if err := json.Unmarshal(datesBytes, &dates); err != nil {
			log.Printf("JSON unmarshal failed: %v", err)
		}
		lesson.Dates = dates

		lessonIds = append(lessonIds, lesson.ID)
		lessons = append(lessons, repository.LessonWithScores{
			Lesson:   lesson,
			Students: make([]repository.StudentWithScores, 0),
		})
	}
	if err := rowsLessons.Err(); err != nil {
		return nil, err
	}

	scores, err := r.getScoresForStudent(ctx, groupId, lessonIds)

	for i, lesson := range lessons {
		for _, student := range students {
			s := repository.StudentWithScores{
				Student:  student,
				Criteria: make(map[string]int),
			}

			if score, ok := scores[lesson.ID][student.ID]; ok {
				s.Criteria = score.Criteria
				s.Progress = score.Progress
			}

			lesson.Students = append(lesson.Students, s)
		}
		lessons[i] = lesson
	}

	return lessons, nil
}

type LessonScoresByStudent map[int]map[int]AggCriteria

type AggCriteria struct {
	Criteria map[string]int
	Progress int
}

func (r *postgresRepo) getScoresForStudent(ctx context.Context, groupId int, lessonIds []int) (LessonScoresByStudent, error) {
	if len(lessonIds) == 0 {
		return nil, nil
	}

	args := make([]interface{}, len(lessonIds)+2)
	args[0] = groupId
	args[1] = groupId

	placeholders := make([]string, len(lessonIds))
	for i, v := range lessonIds {
		placeholders[i] = fmt.Sprintf("$%d", i+3)
		args[i+2] = v
	}

	query := fmt.Sprintf(`
		WITH criteria_max_score AS (
			SELECT g.group_id, SUM(c.max_score) AS total_max_score
			FROM group_criteria g
			JOIN criteria c ON g.criteria_id = c.id
			WHERE g.group_id = $1
			GROUP BY g.group_id
		),
		student_progress AS (
			SELECT
				ls.lesson_id,
				ls.student_id,
				JSONB_OBJECT_AGG(ls.criteria_id::TEXT, ls.score) AS criteria,
				SUM(ls.score) AS total_score,
				l.group_id
			FROM lesson_scores ls
			JOIN lessons l ON ls.lesson_id = l.id and l.group_id = $2
			GROUP BY l.group_id, ls.lesson_id, ls.student_id
		)
		SELECT
			sp.lesson_id,
			sp.student_id,
			sp.criteria,
			(sp.total_score * 100 / cm.total_max_score)::INT AS progress
		FROM student_progress sp
		JOIN criteria_max_score cm ON sp.group_id = cm.group_id
		WHERE sp.lesson_id IN (%s);
    `, strings.Join(placeholders, ", "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		fmt.Printf("scores: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	details := make(LessonScoresByStudent)
	for rows.Next() {
		var criteriaBytes []byte
		var score struct {
			LessonID  int
			StudentID int
			Progress  int
		}
		if err := rows.Scan(&score.LessonID, &score.StudentID, &criteriaBytes, &score.Progress); err != nil {
			return nil, err
		}

		criteria := make(map[string]int)
		if err := json.Unmarshal(criteriaBytes, &criteria); err != nil {
			log.Printf("JSON unmarshal failed: %v", err)
		}

		if _, ok := details[score.LessonID]; !ok {
			details[score.LessonID] = make(map[int]AggCriteria)
		}

		details[score.LessonID][score.StudentID] = AggCriteria{
			Criteria: criteria,
			Progress: score.Progress,
		}
	}

	return details, rows.Err()
}

func (r *postgresRepo) CreateLesson(ctx context.Context, groupId int) (int, error) {
	var id int

	query := `INSERT INTO lessons (group_id) VALUES ($1) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, groupId).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreateLesson exec error: %w", err)
	}
	return id, nil
}

func (r *postgresRepo) DeleteLesson(ctx context.Context, groupId, lessonId int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM lessons WHERE group_id = $1 AND id = $2`, groupId, lessonId)
	if err != nil {
		return fmt.Errorf("DeleteLesson exec error: %w", err)
	}

	return nil
}

func (r *postgresRepo) ChangeLessonDates(ctx context.Context, groupId, lessonId int, dates []string) error {
	datesJSON, err := json.Marshal(dates)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `UPDATE lessons SET dates = $1 WHERE group_id = $2 AND id = $3`, datesJSON, groupId, lessonId)
	if err != nil {
		return fmt.Errorf("ChangeLessonDates exec error: %w", err)
	}

	return nil
}

func (r *postgresRepo) UpdateLessonScore(ctx context.Context, groupId, lessonId, studentId, criteriaId, score int) (int, error) {
	query := `
        INSERT INTO lesson_scores (lesson_id, student_id, criteria_id, score)
             VALUES ($1, $2, $3, $4)
             ON CONFLICT (lesson_id, student_id, criteria_id)
             DO UPDATE SET score = EXCLUDED.score
    `
	_, err := r.db.ExecContext(ctx, query, lessonId, studentId, criteriaId, score)
	if err != nil {
		return 0, fmt.Errorf("UpdateLessonScore exec error: %w", err)
	}

	return r.getScoresByStudent(ctx, groupId, lessonId, studentId)
}

func (r *postgresRepo) DeleteLessonScore(ctx context.Context, groupId, lessonId, studentId, criteriaId int) (int, error) {
	_, err := r.db.ExecContext(ctx, `DELETE FROM lesson_scores WHERE lesson_id = $1 AND student_id = $2 AND criteria_id = $3`, lessonId, studentId, criteriaId)
	if err != nil {
		return 0, fmt.Errorf("DeleteLessonScore exec error: %w", err)
	}

	return r.getScoresByStudent(ctx, groupId, lessonId, studentId)
}

func (r *postgresRepo) getScoresByStudent(ctx context.Context, groupId, lessonId, studentId int) (int, error) {
	rows, err := r.db.QueryContext(ctx, `
		WITH criteria_max_score AS (
			SELECT SUM(c.max_score) AS total_max_score
			FROM group_criteria g
			JOIN criteria c ON g.criteria_id = c.id
			WHERE g.group_id = $1
		),
		student_progress AS (
			SELECT
				SUM(ls.score) AS total_score
			FROM lesson_scores ls
			WHERE ls.lesson_id = $2 AND ls.student_id = $3
			GROUP BY ls.student_id
		)
		SELECT
			(sp.total_score * 100 / cm.total_max_score)::INT AS progress
		FROM student_progress sp
		CROSS JOIN criteria_max_score cm;
 `, groupId, lessonId, studentId)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var progress int
	for rows.Next() {
		if err := rows.Scan(&progress); err != nil {
			return 0, err
		}
	}
	return progress, rows.Err()
}

func (r *postgresRepo) GetRatings(ctx context.Context) ([]entity.Rating, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT student_name, group_name, month, performance_percentage FROM ratings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]entity.Rating, 0, 100)
	for rows.Next() {
		var r entity.Rating
		if err := rows.Scan(&r.Name, &r.GroupName, &r.Month, &r.Performance); err != nil {
			return nil, err
		}
		users = append(users, r)
	}
	return users, rows.Err()
}
