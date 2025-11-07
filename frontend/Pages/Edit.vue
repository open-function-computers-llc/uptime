<template>
<div class="container mt-3">

    <Head title="Edit URL" />

    <div class="row">
        <div class="col">
            <h1>Edit </h1>

            <div class="w-50">
                <TextInput v-model="form.url" label="URL" />
                <TextareaInput v-model="form.meta" label="Meta" />
            </div>

            <div class="d-flex gap-1">
                <button @click="submit" class="btn btn-success" :disabled="isDisabled">Save</button>
                <Link href="/" class="btn btn-danger">Cancel</Link>
            </div>
        </div>
    </div>
</div>
</template>

<script setup>
import { Head, Link, useForm } from "@inertiajs/vue3";
import Layout from "@/Layouts/Standard.vue";
import TextInput from "@/Components/TextInput.vue";
import TextareaInput from "@/Components/TextareaInput.vue";
import { computed } from "vue";

defineOptions({ layout: Layout });

const props = defineProps({
    site: {
        type: Object,
        default: {},
    },
});

const form = useForm({
    url: props.site.URL,
    meta: props.site.Meta,
});

const submit = () => {
    form.post("/update/" + props.site.ID);
}

const isDisabled = computed(() => {
    if (form.url.length > 0 && form.meta.length > 0) {
        return false;
    }

    return true;
});
</script>
