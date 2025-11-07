<template>
<div>
    <label v-if="label" class="form-label" :for="id">
        {{ label }}
        <template v-if="required">*</template>
    </label>
    <select :style="'width:' + width" class="form-select" :value="modelValue" :required="required" :disabled="disabled" @input="$emit('update:modelValue', $event.target.value)" ref="input">
        <option value="" disabled v-html="placeholder"></option>
        <option v-for="(option, i) in options" :value="parseOptionValue(option)" :key="'option-' + i">
            {{ parseOptionText(option) }}
        </option>
    </select>
</div>
</template>

<script setup>
import { ref } from 'vue';


const props = defineProps({
    modelValue: {
        type: [String, Number],
        default: ""
    },
    options: {
        type: [Array, Object],
        default: [],
    },
    optionText: {
        type: String,
        default: "text?",
    },
    optionTextPrepend: {
        type: String,
        default: "",
    },
    optionValue: {
        type: String,
        default: "",
    },
    label: {
        type: String,
        default: "",
    },
    id: {
        type: String,
        default: 'select-input'
    },
    required: {
        type: Boolean,
        default: false
    },
    autofocus: {
        type: Boolean,
        default: false,
    },
    placeholder: {
        type: String,
        default: "Please Choose:",
    },
    width: {
        type: String,
        default: '100px'
    },
    disabled: {
        type: Boolean,
        default: false
    },
});

const input = ref(null);
const parseOptionText = (o) => {
    if (o[props.optionText]) {
        return props.optionTextPrepend + o[props.optionText];
    }

    return o;
};
const parseOptionValue = (o) => {
    if (o[props.optionValue]) {
        return o[props.optionValue];
    }

    return o;
};
</script>

<style lang="scss" scoped>
@import "../scss/variables.scss";

input {
    border-radius: 0;
    border-color: $c-lightGray;
}

label {
    font-weight: bold;
    margin-bottom: 0.25rem;
}
</style>
