<script setup lang="ts">
import { ref } from 'vue';
import { FilterMatchMode } from '@primevue/core/api';

const { $config } = useNuxtApp();

const filters = ref({
  name: { value: null, matchMode: FilterMatchMode.CONTAINS },
  group: { value: null, matchMode: FilterMatchMode.IN },
});

const { data: groups } = await useFetch<{ data: string[] }>(`${$config.public.API_URL}/groups`, {
  transform: (response: { data: string[] }) => response.data.map(group => group.name),
});

const { data: students } = await useFetch<{ data: StudentPerformance[] }>(`${$config.public.API_URL}/ratings`, {
  transform: (response: { data: StudentPerformance[] }) => response.data.map(student => ({
    name: student['student_name'],
    group: student['group_name'],
    type: student['type'],
    performance: `${student['performance_percentage']}%`,
  }) as StudentPerformance),
});
</script>

<template>
  <div class="card">
    <DataTable
        v-model:filters="filters"
        :value="students"
        :rows="15"
        dataKey="id"
        filterDisplay="row"
        paginator
    >
      <template #header>
        <div class="flex flex-wrap gap-2 items-center justify-between">
          <h2 class="m-0">🏆 Rating</h2>
        </div>
      </template>
      <template #empty> No students found.</template>
      <template #loading> Loading students data. Please wait.</template>
      <Column :showFilterMenu="false" field="name" header="Name" style="max-width: 4rem">
        <template #filter="{ filterModel, filterCallback }">
          <InputText v-model="filterModel.value" type="text" @input="filterCallback()" placeholder="Search by name" />
        </template>
      </Column>
      <Column :showFilterMenu="false" field="group" header="Group" style="max-width: 4rem">
        <template #filter="{ filterModel, filterCallback }">
          <MultiSelect
              v-model="filterModel.value"
              :options="groups"
              :max-selected-labels="2"
              @change="filterCallback()"
              placeholder="Any"
          />
        </template>
      </Column>
      <Column field="performance" header="January"></Column>
    </DataTable>
  </div>
</template>


