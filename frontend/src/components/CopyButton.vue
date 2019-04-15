<template>
    <button class="button is-success clipboard" v-bind:data-clipboard-text="url">{{ text }}</button>
</template>

<script>
    export default {
        name: "CopyButton",
        props: ['text', "url"],
    }
</script>

<style scoped>

</style>
