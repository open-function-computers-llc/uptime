<template>
<div class="success-button-wrapper">
    <!-- Success Alert -->
    <div v-if="success" class="alert alert-success alert-dismissible fade show" role="alert">
        <strong>Success!</strong> {{ successMessage }}
        <button type="button" class="btn-close" @click="closeSuccess" aria-label="Close"></button>
    </div>

    <!-- Error Alert -->
    <div v-if="error" class="alert alert-danger alert-dismissible fade show" role="alert">
        <strong>Error!</strong> {{ error }}
        <button type="button" class="btn-close" @click="error = null" aria-label="Close"></button>
    </div>

    <!-- Action Button -->
    <button
        class="btn btn-warning text-white"
        :disabled="loading"
        @click="triggerCleanup">
        <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
        {{ loading ? 'Processing...' : 'Clean Up Outages' }}
    </button>
</div>
</template>

<script setup>
import { ref } from 'vue';

const loading = ref(false);
const success = ref(false);
const successMessage = ref('');
const error = ref(null);

const closeSuccess = () => {
    success.value = false;
    successMessage.value = '';
};

const triggerCleanup = async () => {
    success.value = false;
    successMessage.value = '';
    error.value = null;
    loading.value = true;

    try {
        const csrfToken = document.querySelector('meta[name="csrf-token"]')?.getAttribute('content');

        const response = await fetch('/cleanup', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-CSRF-TOKEN': csrfToken, // Crucial for Laravel/Inertia
            },
        });

        // Check for HTTP errors (4xx, 5xx)
        if (!response.ok) {
            let errMsg = `Server error ${response.status}`;
            try {
                const data = await response.json();
                errMsg = data.message || errMsg;
            } catch (e) {
                // If body isn't JSON, use status text
                errMsg = response.statusText || errMsg;
            }
            throw new Error(errMsg);
        }

        const data = await response.json();

        success.value = true;
        successMessage.value = data.message || 'Cleanup completed successfully.';

        // Auto-hide success after 5 seconds
        setTimeout(() => {
            success.value = false;
        }, 5000);

    } catch (err) {
        console.error('Cleanup Error:', err);
        error.value = err.message || 'An unexpected error occurred. Check console for details.';
    } finally {
        loading.value = false;
    }
};
</script>

<style scoped>
.success-button-wrapper {
    position: relative;
}

.alert {
    position: absolute;
    top: calc(100% + 0.2rem);
    z-index: 1;
    white-space: nowrap;
}
</style>
