<template>
    <button class="button is-success clipboard" :data-clipboard-text="url">{{ text }}</button>
</template>

<script>
    export default {
        name: "CopyButton",
        props: ['text', "url"],
    }
</script>

<style scoped>
.clipboard {
    margin-top: 3px;
}
</style>
