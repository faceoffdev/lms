<script setup lang="ts">
import {ref} from 'vue';
import {FilterMatchMode} from '@primevue/core/api';

const {$config} = useNuxtApp();

const filters = ref({
  name: {value: null, matchMode: FilterMatchMode.CONTAINS},
  group: {value: null, matchMode: FilterMatchMode.IN},
});

const {data: groups} = await useFetch<{ data: string[] }>(`${$config.public.API_URL}/groups`, {
  transform: (response: { data: string[] }) => response.data.map(group => group.name),
});

const {data: students} = await useFetch<{ data: StudentPerformance[] }>(`${$config.public.API_URL}/ratings`, {
  transform: (response: { data: StudentPerformance[] }) => response.data.map(student => ({
    name: student['student_name'],
    group: student['group_name'],
    type: student['type'],
    performance: `${student['performance_percentage']}%`,
  }) as StudentPerformance),
});
</script>

<template>
  <div class="flex h-screen space-x-1">
    <div class="flex flex-col items-center min-w-16 h-full overflow-hidden text-gray-400">
      <div class="flex items-center justify-center mt-3">
        <i class="pi pi-graduation-cap" style="font-size: 2rem"/>
      </div>

      <div class="flex flex-col items-center mt-3 border-t border-gray-700">
        <Button as="router-link" icon="pi pi-book" variant="text" severity="secondary" class="mt-3" to="/groups/1"/>
        <Button icon="pi pi-chart-bar" variant="text" class="mt-3"/>
      </div>
    </div>

    <card class="w-screen">
      <template #title>🏆 Rating</template>
      <template #content>
        <DataTable
            v-model:filters="filters"
            :value="students"
            :rows="15"
            dataKey="id"
            filterDisplay="row"
            paginator
        >
          <template #empty> No students found.</template>
          <template #loading> Loading students data. Please wait.</template>
          <Column :showFilterMenu="false" field="name" header="Name" style="max-width: 4rem">
            <template #filter="{ filterModel, filterCallback }">
              <InputText v-model="filterModel.value" type="text" @input="filterCallback()"
                         placeholder="Search by name"/>
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
      </template>
    </card>
  </div>
</template>


