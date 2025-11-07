<template><span>{{ d.toLocaleString() }}</span>
</template>

<script setup>
const props = defineProps({
    ts: String,
});

const d = new Date(props.ts);
</script>
