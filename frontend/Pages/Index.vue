<template>
<div class="my-3">

    <Head title="Uptime Monitoring" />

    <div class="d-flex gap-3 justify-content-between align-items-center mb-3">
        <Link class="btn btn-primary text-white" href="/add">
        Add Site
        <PlusLg class="nudge-up" />
        </Link>

        <Switch v-model="showDeleted" :label="showDeleted ? 'Deleted Sites' : 'Active Sites'" />
    </div>

    <p v-if="sitesToShow.length === 0">No sites</p>

    <div class="d-flex flex-wrap gap-3">
        <Site v-for="s in sitesToShow" :site="s" :key="s.ID" @set="(e) => activeSite = e" />
    </div>

    <Modal v-if="activeSite" @close="activeSite = null">
        <SiteDetails :site="activeSite" />
    </Modal>
</div>
</template>

<script setup>
import { Head, Link } from "@inertiajs/vue3";
import Layout from "@/Layouts/Standard.vue";
import Site from "@/Components/Site.vue";
import { computed, onMounted, onUnmounted, ref } from "vue";
import Switch from "@/Components/Switch.vue";
import PlusLg from "@/Icons/PlusLg.vue";
import Modal from "../Components/Modal.vue";
import SiteDetails from "../Components/SiteDetails.vue";

defineOptions({ layout: Layout });

const props = defineProps({
    sites: {
        type: Object,
        default: {},
    },
});

const loadedSites = ref(props.sites);

let intervalId = null;
const fetchSites = async () => {
    try {
        const response = await fetch("/api/load-sites", {
            headers: { "Accept": "application/json" },
        });
        if (!response.ok) {
            throw new Error(`HTTP error! ${response.status}`);
        }
        const data = await response.json();
        loadedSites.value = data;
    } catch (err) {
        console.error("Error fetching sites:", err);
    }
};

const pollSites = () => {
    fetchSites();
    intervalId = setInterval(fetchSites, 15 * 1000); // do it every 15 seconds
};

onMounted(() => {
    pollSites();
});
onUnmounted(() => {
    if (intervalId) clearInterval(intervalId);
});

const showDeleted = ref(false);
const sitesToShow = computed(() => {
    const output = [];

    for (const [siteID, site] of Object.entries(loadedSites.value)) {
        if (showDeleted.value) {
            if (site.IsDeleted) {
                output.push(site);
            }
        } else {
            if (!site.IsDeleted) {
                output.push(site);
            }
        }
    }

    return output.sort((a, b) => {
        // First: prioritize down items (false before true)
        if (a.IsUp !== b.IsUp) {
            return a.IsUp ? 1 : -1;
        }

        const aNormalized = a.URL.replace(/^https?:\/\//, "")
            .replace(/^www\./, "")
            .toLowerCase();
        const bNormalized = b.URL.replace(/^https?:\/\//, "")
            .replace(/^www\./, "")
            .toLowerCase();

        return aNormalized.localeCompare(bNormalized);
    });
});

const activeSite = ref(null);
</script>

<style lang="scss" scoped>
@use "@/scss/variables" as *;

a {
    text-decoration: none;
    color: $c-gray;
}

.nudge-up {
    transform: translateY(-2px);
}
</style>
