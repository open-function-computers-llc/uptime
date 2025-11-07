<template>
<div class="container mt-3">

    <Head title="Purge URL" />

    <div class="row">
        <div class="col">
            <h1>Are you sure?</h1>

            <p class="w-50">There's no going back from this. This will remove the site, it's webhooks, and any outage history. You've been warned!
            </p>

            <div class="d-flex gap-1">
                <button @click="submit" class="btn btn-success">Purge</button>
                <Link href="/" class="btn btn-danger">Cancel</Link>
            </div>
        </div>
    </div>
</div>
</template>

<script setup>
import { Head, Link, useForm } from "@inertiajs/vue3";
import Layout from "@/Layouts/Standard.vue";

defineOptions({ layout: Layout });

const props = defineProps({
    site: {
        type: Object,
        default: {},
    },
});

const submit = () => {
    const form = useForm({});
    form.post("/purge/" + props.site.ID);
}
</script>
