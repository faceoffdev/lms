<script setup lang="ts">
import type { PageState, RatingChangeEvent } from "primevue";
import { range } from "@antfu/utils";

const {$config} = useNuxtApp();
const route = useRoute();

const readonly = !('admin' in route.query);

const tab = { id: Number(route.params.group) } as Group;
const groups = ref<Group[]>([]);
const criteria = ref<Criteria[]>([]);
const students = ref<Student[]>([]);
const lessons = ref<Lesson[]>([]);

const loading = ref(false);

const loadData = async () => {
  loading.value = true;
  try {
    const [groupsResponse, criteriaResponse, studentsResponse, lessonsResponse] = await Promise.all([
      useFetch<{ data: Group[] }>(`${$config.public.API_URL}/groups`),
      useFetch<{ data: Criteria[] }>(`${$config.public.API_URL}/groups/${tab.id}/criteria`),
      useFetch<{ data: Student[] }>(`${$config.public.API_URL}/groups/${tab.id}/students`),
      useFetch<{ data: Lesson[] }>(`${$config.public.API_URL}/groups/${tab.id}/lessons`, {
        transform: (response: { data: Lesson[] }) => response.data.map((lesson) => {
          lesson.dates = lesson.dates.map(d => new Date(d));
          return lesson;
        }),
      }),
    ]);

    groups.value = groupsResponse.data.value['data'];
    criteria.value = criteriaResponse.data.value['data'];
    students.value = studentsResponse.data.value['data'];
    lessons.value = lessonsResponse.data.value;
  } catch (error) {
    console.error('Ошибка загрузки данных:', error);
  } finally {
    loading.value = false;
  }
};

const addLesson = async () => {
  try {
    const {data} = await $fetch<{
      data: { 'lesson_id': number }
    }>(`${$config.public.API_URL}/groups/${tab.id}/lessons`, {method: 'POST'});
    lessons.value.unshift({
      id: Number(data['lesson_id']),
      dates: [],
      students: students.value.map(s => ({
        id: s.id,
        name: s.name,
        progress: 0,
        criteria: {}
      })),
    });
  } catch (error) {
    console.error("Ошибка при добавлении урока:", error);
  }
};

const delLesson = async (lessonId: number) => {
  try {
    await $fetch(`${$config.public.API_URL}/groups/${tab.id}/lessons/${lessonId}`, {method: 'DELETE'});

    lessons.value = lessons.value.filter((lesson) => lesson.id !== lessonId);
  } catch (error) {
    console.error("Ошибка при удалении урока:", error);
  }
};

const changeCriteria = async (event: RatingChangeEvent, studentId: number, lessonId: number, critId: number) => {
  let method = 'DELETE';
  const url = `${$config.public.API_URL}/groups/${tab.id}/lessons/${lessonId}/criteria`;
  const body = {
    'student_id': studentId,
    'criteria_id': critId
  };

  if (event.value && event.value > 0) {
    method = 'PATCH';
    body['score'] = event.value;
  }

  try {
    const {data} = await $fetch<{ data: { progress: number } }>(url, {method, body: JSON.stringify(body)});

    const lesson = lessons.value.find(l => l.id === lessonId);

    if (lesson) {
      const student = lesson.students.find(s => s.id === studentId);

      if (student) {
        student.progress = data.progress;
      }
    }
  } catch (error) {
    console.error("Ошибка при изменении данных:", error);
  }
};

const onDateSelect = async (lessonId, dates: Date[]) => {
  dates = dates?.map(d => d.toISOString()) ?? [];

  try {
    await $fetch(`${$config.public.API_URL}/groups/${tab.id}/lessons/${lessonId}/dates`, {
      method: 'PATCH',
      body: JSON.stringify({dates})
    });
  } catch (error) {
    console.error("Ошибка при изменении дат:", error);
  }
};

const navigateToByGroup = (id: number): string => {
  let url = `/groups/${id}`;

  if ('admin' in route.query) {
    url += `?admin`;
  }

  return navigateTo(url);
}

const handlePageChange = (event: PageState) => {
  currentPage.value = event.page + 1;
};

loadData();

const currentPage = ref(1);
const lessonsPerPage = 3;
const paginatedLessons = computed(() => {
  const start = (currentPage.value - 1) * lessonsPerPage;
  const end = start + lessonsPerPage;

  return lessons.value.slice(start, end);
});
</script>

<template>
  <div class="flex w-screen h-screen space-x-1">
    <div class="flex flex-col items-center min-w-16 h-full overflow-hidden text-gray-400">
      <div class="flex items-center justify-center mt-3">
        <i class="pi pi-graduation-cap" style="font-size: 2rem; color: #298a67"/>
      </div>

      <div class="flex flex-col items-center mt-3 border-t border-gray-700">
        <Button icon="pi pi-book" variant="text" class="mt-3"/>
        <Button as="router-link" icon="pi pi-chart-bar" variant="text" severity="secondary" class="mt-3" to="/ratings" />
      </div>
    </div>
    <div class="main-wrapper">
      <Skeleton v-if="loading" width="100%" height="3.5rem"/>
      <DataTable
          v-if="loading"
          v-for="j in range(4)"
          :key="j"
          :value="range(5)"
          :style="{ marginBottom: '25px' }"
          stripedRows
      >
        <template #header>
          <Skeleton width="16rem" height="4rem"/>
        </template>
        <Column v-for="i in range(7)" :key="i">
          <template #header>
            <Skeleton height="3rem"/>
          </template>
          <template #body>
            <Skeleton/>
          </template>
        </Column>
      </DataTable>

      <Tabs v-if="groups.length" :value="tab.id" scrollable>
        <TabList>
          <Tab v-for="group in groups" :key="group.name" :value="group.id" @click="navigateToByGroup(group.id)">
            {{ group.name }}
          </Tab>
        </TabList>
      </Tabs>

      <div v-if="groups.length > 0 && lessons.length === 0 && !readonly" class="button-container">
        <Button @click="addLesson" severity="success" label="Add first lesson"/>
      </div>

      <DataTable
          v-for="(lesson, index) in paginatedLessons"
          :key="lesson.id"
          :value="lesson.students"
          :style="{ marginBottom: '25px' }"
          stripedRows
      >
        <template #header>
          <div class="flex flex-wrap items-center justify-between gap-2"
               style="display: flex; justify-content: space-between; align-items: center;">
            <FloatLabel variant="in">
              <DatePicker
                  v-model="lesson.dates"
                  :manualInput="false"
                  :disabled="readonly"
                  @value-change="onDateSelect(lesson.id, $event)"
                  selectionMode="multiple"
                  inputId="in_label"
                  iconDisplay="input"
                  variant="filled"
                  dateFormat="d M"
                  showButtonBar
                  showIcon
              />
              <label for="in_label">Date</label>
            </FloatLabel>
            <Button v-if="index === 0 && currentPage === 1 && !readonly" @click="addLesson" severity="success" icon="pi pi-plus"/>
            <Button v-else-if="!readonly" @click="delLesson(lesson.id)" icon="pi pi-trash" severity="danger"/>
          </div>
        </template>
        <Column field="name" header="ФИО" style="width: 15rem" sortable></Column>
        <Column v-for="crit in criteria" :key="crit.id" :field="`crit_${crit.id}`" :header="crit.name">
          <template #body="slotProps">
            <Rating
                v-model="slotProps.data.criteria[crit.id]"
                :readonly="readonly"
                :stars="crit.max_score"
                @change="changeCriteria($event, slotProps.data.id, lesson.id, crit.id)"
            />
          </template>
        </Column>
        <Column field="progress" header="% успеха" sortable></Column>
      </DataTable>
      <Paginator
          v-if="lessons.length > lessonsPerPage"
          :rows="lessonsPerPage"
          :totalRecords="lessons.length"
          @page="handlePageChange"
      />
    </div>
  </div>
</template>

<style scoped>
.button-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>