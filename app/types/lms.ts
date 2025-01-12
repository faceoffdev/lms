interface Criteria {
    id: number;
    name: string;
    max_score: number;
}

interface Group {
    id: number;
    name: string;
}

interface Lesson {
    id: number;
    dates: Date[];
    students: Student[];
}

interface Student {
    id: number;
    name: string;
    progress: number;
    criteria: Record<number, number>;
}
