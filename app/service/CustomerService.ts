export const CustomerService = {

    getData() {
        return [
            {
                name: "Даня Матушкин",
                group: "B1+",
                type: "group",
                performance: 100
            },
            {
                name: "Алина Гасанова",
                group: "B1+",
                type: "group",
                performance: 96.43
            },
            {
                name: "Виктор Коробков",
                group: "B1+",
                type: "group",
                performance: 91.43
            },
            {
                name: "Николай Шестаков",
                group: "B1+",
                type: "group",
                performance: 67.86
            }
        ];
    },

    getCustomersMedium() {
        return Promise.resolve(this.getData().slice(0, 50));
    },
};
