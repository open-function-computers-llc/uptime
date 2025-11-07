<template>
<main>
    <header>
        <div class="container-fluid">
            <div class="row">
                <div class="col py-3">
                    <div class="logo-wrapper">
                        <OfcLogo height="24" width="24" />
                    </div>
                    Uptime Monitor
                </div>
            </div>
        </div>
    </header>
    <div class="container-fluid">
        <div class="row">
            <div class="col">

                <slot />

            </div>
        </div>
    </div>
</main>
</template>

<script setup>
import OfcLogo from "@/Components/OfcLogo.vue";
</script>

<style lang="scss" scoped>
@use "@/scss/variables" as *;

header {
    background-color: $c-gray;
    color: $c-white;
    font-weight: bold;
}

.logo-wrapper {
    background-color: $c-white;
    display: inline-block;
    padding: 0.5rem;
    border-radius: 50%;
    margin-right: 1rem;
}
</style>
