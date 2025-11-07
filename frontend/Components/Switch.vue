<template>
<div class="form-check form-switch d-flex flex-row-reverse gap-2">
    <input class="form-check-input" type="checkbox" role="switch" id="flexSwitchCheckDefault" v-model="inputValue">
    <label class="form-check-label d-block me-5" for="flexSwitchCheckDefault">
        {{ label }}
    </label>
</div>
</template>

<script setup>
import { ref, watch } from 'vue';

const props = defineProps({
    label: {
        type: String,
        default: ""
    },
    modelValue: {
        type: Boolean,
        default: false,
    },
});
const emit = defineEmits(["update:modelValue"]);
const inputValue = ref(props.modelValue);

watch(inputValue, () => {
    emit('update:modelValue', inputValue.value);
});
</script>
