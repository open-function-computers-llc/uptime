<template>
<div>
    <div class="container-fluid">
        <div class="row">
            <div class="col">
                <h3 class="text-center">{{ site.URL }}</h3>
                <div v-html="nl2br(site.Meta)" />
                <ul>
                    <li :class="{ warning: site.Uptime_1day < 99, danger: site.Uptime_1day < 98 }">
                        <span>1 Day:</span> {{ site.Uptime_1day }}%
                    </li>
                    <li :class="{ warning: site.Uptime_7day < 99, danger: site.Uptime_7day < 98 }">
                        <span>7 Day:</span> {{ site.Uptime_7day }}%
                    </li>
                    <li :class="{ warning: site.Uptime_30day < 99, danger: site.Uptime_30day < 98 }">
                        <span>30 Day:</span> {{ site.Uptime_30day }}%
                    </li>
                    <li :class="{ warning: site.Uptime_60day < 99, danger: site.Uptime_60day < 98 }">
                        <span>60 Day:</span> {{ site.Uptime_60day }}%
                    </li>
                    <li :class="{ warning: site.Uptime_90day < 99, danger: site.Uptime_90day < 98 }">
                        <span>90 Day:</span> {{ site.Uptime_90day }}%
                    </li>
                </ul>
            </div>
        </div>
    </div>

    <template v-if="!isLoading">
        <div class="container-fluid">
            <div class="row">
                <div class="col">
                    <h4>Outages</h4>

                    <ul class="list-group outages-list">
                        <li v-for="o in siteDetails.outages"
                            class="list-group-item"
                            :class="{
                                'list-group-item-warning': o.duration > 60,
                                'list-group-item-danger': o.duration > 180,
                            }"
                            :key="o.ID">
                            <Timestamp :ts="o.start" />:
                            {{ o.duration }} second(s) {{ o }}
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </template>

    <Spinner v-if="isLoading" />
</div>
</template>

<script setup>
import { ref } from 'vue';
import Spinner from "@/Components/Spinner.vue";
import Timestamp from "@/Components/Timestamp.vue";

const props = defineProps({
    site: {
        type: Object,
        default: {},
    },
});

const siteDetails = ref(null);
const isLoading = ref(false);

const fetchdetails = async (id) => {
    isLoading.value = true;
    try {
        const response = await fetch("/details/" + id, {
            headers: { "Accept": "application/json" },
        });

        if (!response.ok) {
            throw new Error(`HTTP error! ${response.status}`);
        }

        const data = await response.json();
        siteDetails.value = data;
        isLoading.value = false;
    } catch (err) {
        console.error("Error fetching sites:", err);
        isLoading.value = false;
    }
};

const loadDetails = () => {
    fetchdetails(props.site.ID);
}
loadDetails();

const nl2br = (s) => {
    return s.split("\n").join("<br />");
}
</script>

<style lang="scss" scoped>
.outages-list {
    overflow-y: auto;
    max-height: calc(100vh - 305px - 2rem - 5vh);
}
</style>
