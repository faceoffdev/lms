<script setup lang="ts">
import {ref} from 'vue';
import {FilterMatchMode} from '@primevue/core/api';

const {$config} = useNuxtApp();

const filters = ref({
  name: {value: null, matchMode: FilterMatchMode.CONTAINS},
  group: {value: null, matchMode: FilterMatchMode.IN},
});

const months = ['May', 'April', 'March', 'February', 'January'];
const currentMonth = ref<string>(months[0]);

const groups = ref<string>([]);
const students = ref<StudentPerformance[]>([]);

const loadData = async (month: string) => {
  const getPerformance = (performance, type) => performance[type] ? `${performance[type]}%` : 'N/A';

  const studentRatings = [] as StudentPerformance[];
  const groupNames = new Set<string>();

  try {
    const ratingsResponse = await useFetch(`${$config.public.API_URL}/ratings?month=${month}`);

    ratingsResponse.data.value['data'].forEach((rating) => {
      groupNames.add(rating['group_name']);

      studentRatings.push({
        name: rating['student_name'],
        group: rating['group_name'],
        performance: {
          'common': getPerformance(rating['performance'], 'common'),
          'group': getPerformance(rating['performance'], 'group'),
          'club': getPerformance(rating['performance'], 'club'),
        },
      });
    });
  } catch (error) {
    console.error('Ошибка загрузки данных:', error);
  }

  students.value = studentRatings;
  groups.value = [...groupNames];
};

loadData(currentMonth.value);
</script>

<template>
  <div class="flex w-screen h-screen space-x-1">
    <div class="flex flex-col items-center min-w-16 h-full overflow-hidden text-gray-400">
      <div class="flex items-center justify-center mt-3">
        <i class="pi pi-graduation-cap" style="font-size: 2rem; color: #298a67"/>
      </div>

      <div class="flex flex-col items-center mt-3 border-t border-gray-700">
        <Button as="router-link" icon="pi pi-book" variant="text" severity="secondary" class="mt-3" to="/groups/1"/>
        <Button icon="pi pi-chart-bar" variant="text" class="mt-3"/>
      </div>
    </div>
    <card class="w-screen">
      <template #title>🏆 Rating</template>
      <template #content>
        <Tabs v-model:value="currentMonth" scrollable>
          <TabList>
            <Tab v-for="m in months" :key="m" :value="m" @click="loadData(currentMonth)">
              {{ m }}
            </Tab>
          </TabList>
        </Tabs>

        <DataTable
            v-model:filters="filters"
            :value="students"
            :rows="15"
            dataKey="id"
            filterDisplay="row"
            selectionMode="single"
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
          <Column :showFilterMenu="false" field="group" header="Group name" style="max-width: 4rem">
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
          <Column header="Group">
            <template #body="slotProps">
              {{ slotProps.data.performance['group'] }}
            </template>
          </Column>
          <Column header="Club">
            <template #body="slotProps">
              {{ slotProps.data.performance['club'] }}
            </template>
          </Column>
          <Column header="Common">
            <template #body="slotProps">
              {{ slotProps.data.performance['common'] }}
            </template>
          </Column>
        </DataTable>
      </template>
    </card>
  </div>
</template>


