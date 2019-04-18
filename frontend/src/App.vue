<template>
    <section class="section">
        <div class="container">
            <h1 class="title">yitu</h1>

            <p class="subtitle">
                Work in process. Testing only. Delete at any time. Max file size 50 MiB.
                开发中，仅供测试，随时删库，最大文件大小50MB。
            </p>

            <vue-dropzone ref="myVueDropzone" id="dropzone" :options="dropzoneOptions"></vue-dropzone>

        </div>
    </section>
</template>

<script>
    const VERSION = `v1.0.0-dev`;

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
        mounted: function () {
            window.console.log(`yitu ${VERSION}`)
            new ClipboardJS('.clipboard');

            document.onpaste = ((event) => {
                const items = (event.clipboardData || event.originalEvent.clipboardData).items;
                for (const index in items) {
                    const item = items[index];
                    if (item.kind === 'file') {
                        this.$refs.myVueDropzone.addFile(item.getAsFile())
                    }
                }
            })
        },
        data: function () {
            return {
                VERSION: VERSION,
                dropzoneOptions: {
                    paramName: "tu",
                    maxFilesize: 50,
                    url: 'https://t.halu.lu/api/upload',
                    timeout: 0,
                    acceptedFiles: "image/*",
                    success: ((file, response) => {
                        window.console.log(response);
                        const url = response.url;
                        file.dataURL = url;

                        let urlDiv = file.previewElement.querySelector(".dz-filename").cloneNode(true);
                        urlDiv.querySelector("span[data-dz-name]").textContent = url;
                        file.previewElement.querySelector(".dz-details").appendChild(urlDiv);
                        file.previewElement.querySelector(".dz-details").setAttribute("data-clipboard-text", url);

                        let copyButton = new CopyButtonClass({
                            propsData: {
                                text: 'Copy URL',
                                url: url,
                            }
                        });
                        copyButton.$mount();
                        file.previewElement.querySelector(".dz-details").appendChild(copyButton.$el);

                        let webpButton = new CopyButtonClass({
                            propsData: {
                                text: 'Copy WebP URL',
                                url: url + "/webp",
                            }
                        });
                        webpButton.$mount();
                        file.previewElement.querySelector(".dz-details").appendChild(webpButton.$el);
                    }),
                }
            }
        }
    }
</script>

<style lang="css">
    @import '../node_modules/bulma/css/bulma.css';
</style>
