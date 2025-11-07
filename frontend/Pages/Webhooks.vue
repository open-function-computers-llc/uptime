<template>
<div class="container mt-3">

    <Head title="Edit URL Webhooks" />

    <div class="row">
        <div class="col d-flex gap-3 justify-content-between align-items-center">
            <h1>Webhooks:</h1>

            <Link href="/" class="btn btn-primary">
            <ArrowLeft /> Back</Link>
        </div>
    </div>

    <template v-if="!isAddingWebhook">
        <div class="row">
            <div class="col">
                <p v-if="webhooks.length === 0" class="text-secondary">No webhooks for {{ site.URL }}</p>
                <ul v-else class="list-group list-group-flush">
                    <li v-for="hook in webhooks" class="list-group-item list-group-item-action d-flex align-items-center justify-content-between gap-1" :key="hook.ID">
                        <span>
                            <span class="fw-bold">{{ hook.Name }}</span>: {{ hook.URL }}
                        </span>

                        <button @click="deleteWebhook(hook.ID)" class="btn btn-danger btn-sm">
                            <Trash />
                        </button>
                    </li>
                </ul>

                <button @click="isAddingWebhook = true" class="btn btn-success text-white">
                    Add Webhook
                    <PlusLg class="nudge-up" />
                </button>
            </div>
        </div>

    </template>

    <template v-else>
        <div class="row">
            <div class="col">
                <TextInput v-model="form.name" label="Webhook Name" />
            </div>

            <div class="col">
                <TextInput v-model="form.url" label="URL" />
            </div>

            <div class="col-auto">
                <SelectInput v-model="form.verb"
                    :options="['GET', 'POST']"
                    label="Verb" />
            </div>

            <div class="col-auto">
                <SelectInput v-model="form.hookType"
                    :options="['Standard', 'Emergency']"
                    label="Type" />
            </div>
        </div>

        <div class="row">

            <div class="d-flex gap-1">
                <button @click="storeWebhook" class="btn btn-success text-white">Save</button>
                <button @click="isAddingWebhook = false" class="btn btn-danger">Cancel</button>
            </div>
        </div>
    </template>

</div>
</template>

<script setup>
import { Head, Link, useForm } from "@inertiajs/vue3";
import Layout from "@/Layouts/Standard.vue";
import TextInput from "@/Components/TextInput.vue";
import SelectInput from "@/Components/SelectInput.vue";
import { ref } from "vue";
import PlusLg from "@/Icons/PlusLg.vue";
import Trash from "@/Icons/Trash.vue";
import ArrowLeft from "@/Icons/ArrowLeft.vue";

defineOptions({ layout: Layout });

const props = defineProps({
    webhooks: {
        type: Array,
        default: [],
    },
    site: {
        type: Object,
        default: {},
    },
});

const isAddingWebhook = ref(false);
const form = useForm({
    siteID: props.site.ID,
    name: "",
    url: "",
    verb: "POST",
    hookType: 1,
});


const storeWebhook = () => {
    form.post("/store-webhook/" + props.site.ID, {
        preserveState: false,
    });
}

const deleteWebhook = (id) => {
    const deleteForm = useForm({
        siteID: props.site.ID,
        webhookID: id,
    });
    deleteForm.delete("/delete-webhook", {
        preserveState: false,
    });
}
</script>
