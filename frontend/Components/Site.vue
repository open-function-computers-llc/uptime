<template>
<div class="site" :class="backgroundColor">
    <div class="fw-bold">
        {{ site.URL }}
    </div>
    <div class="d-flex justify-content-between align-items-end">
        <ul>
            <li :class="{ 'text-warning fw-bold': site.Uptime_1day < 99, 'text-danger': site.Uptime_1day < 98 }">
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
        <div class="d-flex gap-1 flex-column">
            <template v-if="!site.IsDeleted">
                <a class="btn btn-info btn-sm" href="#" @click.prevent="setActive" title="Show Info">
                    <Info />
                </a>
                <Link class="btn btn-primary btn-sm" :href="'/edit/' + site.ID" title="Edit">
                <Pencil />
                </Link>
                <Link class="btn btn-success btn-sm" :href="'/webhooks/' + site.ID" title="Edit Webhooks">
                <CloudUpload />
                </Link>
                <button class="btn btn-danger btn-sm" @click="removeSite" title="Delete">
                    <Trash />
                </button>
            </template>

            <template v-else>
                <button class="btn btn-success btn-sm restore" @click="restoreSite" title="Restore Site">
                    <Trash />
                    <CircleSlash />
                </button>
                <Link class="btn btn-danger btn-sm" :href="'/purge/' + site.ID" title="Purge Site">
                <CircleSlash />
                </Link>
            </template>
        </div>
    </div>

</div>
</template>

<script setup>
import { Link, useForm } from "@inertiajs/vue3";
import { computed } from 'vue';
import Trash from "@/Icons/Trash.vue";
import CircleSlash from "@/Icons/CircleSlash.vue";
import Pencil from "@/Icons/Pencil.vue";
import Info from "@/Icons/Info.vue";
import CloudUpload from "@/Icons/CloudUpload.vue";

const props = defineProps({
    site: {
        type: Object,
        default: {},
    },
});

const backgroundColor = computed(() => {
    const output = [];
    if (props.site.IsUp) {
        output.push("green");
    } else {
        output.push("red");
    }

    if (props.site.IsDeleted) {
        output.push("deleted");
    }

    return output;
});

const removeSite = () => {
    const form = useForm({});
    form.post("/remove/" + props.site.ID, {
        preserveState: false,
    });
};
const restoreSite = () => {
    const form = useForm({});
    form.post("/restore/" + props.site.ID, {
        preserveState: false,
    });
};

const emit = defineEmits(['set']);
function setActive() {
    emit('set', props.site);
}
</script>

<style lang="scss" scoped>
@use "@/scss/variables" as *;

.site {
    aspect-ratio: 1.75;
    background-color: $c-gray;
    width: 340px;
    padding: 1rem;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    border: 2px solid;

    &.green {
        background-color: lighten($c-green, 50);
        border-color: $c-green;
    }

    &.red {
        background-color: lighten($c-red, 50);
        border-color: $c-red;
    }

    &.deleted {
        background-color: lighten($c-lightGray, 10);
        border-color: $c-lightGray;
    }
}

ul {
    margin: 0;
    padding: 0;
    list-style-type: none;

    span {
        display: inline-block;
        width: 60px;
    }
}

.restore {
    position: relative;

    svg {
        &:first-of-type {
            scale: 0.75;
        }

        &:last-of-type {
            position: absolute;
            top: 0.4rem;
            left: 0.5rem;
            scale: 1.3;
        }
    }
}
</style>
