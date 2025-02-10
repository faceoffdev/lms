<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { FilterMatchMode } from '@primevue/core/api';
import { CustomerService } from '@/service/CustomerService';

const { $config } = useNuxtApp();

const filters = ref({
  name: { value: null, matchMode: FilterMatchMode.CONTAINS },
  group: { value: null, matchMode: FilterMatchMode.EQUALS },
});

const students = ref();

const { data: groups } = await useFetch<{ data: string[] }>(`${$config.public.API_URL}/groups`, {
  transform: (response: { data: string[] }) => response.data.map(group => group.name),
});

const loading = ref(true);

onMounted(() => {
  CustomerService.getCustomersMedium().then(data => {
    students.value = data;
    loading.value = false;
  });
});
</script>

<template>
  <div class="card">
    <DataTable
        v-model:filters="filters"
        :value="students"
        :loading="loading"
        :rows="15"
        dataKey="id"
        filterDisplay="row"
        paginator
    >
      <template #header>
        <div class="flex flex-wrap gap-2 items-center justify-between">
          <h4 class="m-0">Rating</h4>
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
          <Select
              v-model="filterModel.value"
              :options="groups"
              @change="filterCallback()"
              placeholder="Any"
          />
        </template>
      </Column>
      <Column field="performance" header="January"></Column>
    </DataTable>
  </div>
</template>


