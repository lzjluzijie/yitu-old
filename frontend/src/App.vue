<template>
    <section class="section">
        <div class="container">
            <h1 class="title">yitu</h1>

            <p class="subtitle">
                Work in progress. Max file size 50 MiB.&nbsp;
                开发中，最大文件大小 50MiB。
            </p>

            <vue-dropzone
                ref="myVueDropzone"
                id="dropzone"
                :options="dropzoneOptions"
                @vdropzone-success="success"
            ></vue-dropzone>

        </div>
    </section>
</template>

<script>
    const VERSION = `v1.0.4`;

    import Vue from 'vue'
    import ClipboardJS from 'clipboard';
    import vue2Dropzone from 'vue2-dropzone'
    import 'vue2-dropzone/dist/vue2Dropzone.min.css'
    import CopyButton from './components/CopyButton.vue'

    const CopyButtonClass = Vue.extend(CopyButton);

    export default {
        name: 'app',
        components: {
            vueDropzone: vue2Dropzone
        },
        mounted() {
            console.log(`yitu ${VERSION}`)
            new ClipboardJS('.clipboard');

            document.onpaste = ((event) => {
                const items = Array.from((event.clipboardData || event.originalEvent.clipboardData).items);
                items.forEach(item => {
                    if (item.kind === 'file') {
                        this.$refs.myVueDropzone.addFile(item.getAsFile())
                    }
                })
            })
        },
        data() {
            return {
                VERSION,
                dropzoneOptions: {
                    paramName: "tu",
                    maxFilesize: 50,
                    url: 'https://t.halu.lu/api/upload',
                    timeout: 0,
                    acceptedFiles: "image/*",
                }
            }
        },
        methods: {
            success(file, response) {
                const url = response.url;
                file.dataURL = url;

                let urlDiv = file.previewElement.querySelector(".dz-filename").cloneNode(true);
                urlDiv.querySelector("span[data-dz-name]").textContent = url;
                const details = file.previewElement.querySelector(".dz-details")
                details.appendChild(urlDiv);
                details.setAttribute("data-clipboard-text", url);

                let copyButton = new CopyButtonClass({
                    propsData: {
                        text: 'Copy URL',
                        url,
                    }
                });
                copyButton.$mount();
                details.appendChild(copyButton.$el);

                let webpButton = new CopyButtonClass({
                    propsData: {
                        text: 'Copy WebP URL',
                        url: url + "/webp",
                    }
                });
                webpButton.$mount();
                details.appendChild(webpButton.$el);
            }
        }
    }
</script>

<style lang="css">
    @import '../node_modules/bulma/css/bulma.css';
</style>
