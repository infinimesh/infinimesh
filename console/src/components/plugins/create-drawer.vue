<template>
    <n-button @click="show = true" type="success" dashed>
        <template #icon>
            <n-icon>
                <add-outline />
            </n-icon>
        </template>
        Register Plugin
    </n-button>
    <n-drawer v-model:show="show" width="80%">
        <n-drawer-content>
            <template #header> Register Plugin </template>
            <template #footer>
                <n-space justify="end" align="center">
                    <n-button type="error" round secondary @click="show = false">Cancel</n-button>
                    <n-button type="info" round secondary @click="reset">Reset</n-button>
                    <n-button type="warning" round @click="handleSubmit">Submit</n-button>
                </n-space>
            </template>
            <n-grid cols="2" responsive="self" :x-gap="20">
                <n-grid-item span="s:2 m:2 l:1 xl:1 2xl:1">
                    <n-form ref="form" :model="model" :rules="rules" label-placement="top">
                        <n-form-item label="Title" path="title">
                            <n-input v-model:value="model.title" placeholder="Make it bright" />
                        </n-form-item>
                        <n-form-item label="Public" path="public">
                            <n-switch v-model:value="model.public">
                                <template #checked>
                                    Will be available to all platform namespaces
                                </template>
                                <template #unchecked>
                                    Will be only available to partiular namespaces
                                </template>
                            </n-switch>
                        </n-form-item>
                        <n-form-item label="Logo" path="logo">
                            <n-input v-model:value="model.logo" placeholder="App/Plugin Logo image URL" />
                        </n-form-item>
                        <n-form-item label="Description" path="description">
                            <n-input v-model:value="model.description" placeholder="Your App/Plugin Description"
                                type="textarea" rows="15" />
                        </n-form-item>
                        <n-alert title="Markdown syntax supported" type="info">
                            <template #icon>
                                <n-icon>
                                    <logo-markdown />
                                </n-icon>
                            </template>
                            Don't make it too long and
                            Start with level 3(###) if you're using headers :)
                        </n-alert>

                        <n-form-item label="Kind" path="kind" style="margin-top: 5px">
                            <n-radio-group v-model:value="model.kind" name="kind">
                                <n-radio-button value="UNKNOWN" :disabled="true" label="Unknown" />
                                <n-radio-button value="EMBEDDED" label="Embedded" />
                            </n-radio-group>
                        </n-form-item>
                        <n-alert title="Note" type="info">
                            <template #icon>
                                <n-icon>
                                    <bookmark-outline />
                                </n-icon>
                            </template>
                            Hover your cursor over kind label in preview to see the differences
                        </n-alert>

                        <n-form-item label="Frame URL" path="embedded_conf.frame" style="margin-top: 5px">
                            <n-input v-model:value="model.embedded_conf.frame" placeholder="IFrame URL to embed" />
                        </n-form-item>
                    </n-form>

                    <n-alert title="Error creating App/Plugin" type="error" v-if="error">
                        {{ error }}
                    </n-alert>
                </n-grid-item>
                <n-grid-item span="s:2 m:2 l:1 xl:1 2xl:1">
                    <n-grid cols="24">
                        <n-grid-item span="1">
                            <n-divider vertical dashed style="height: 100%;" />
                        </n-grid-item>
                        <n-grid-item span="23">
                            <h3>Marketplace Preview:</h3>
                            <plugin-card :plugin="model" :preview="true" />
                        </n-grid-item>
                    </n-grid>
                </n-grid-item>
            </n-grid>

        </n-drawer-content>
    </n-drawer>
</template>

<script setup>
import { ref, watch } from "vue"
import {
    useLoadingBar, NButton, NIcon, NDrawer, NDrawerContent, NSpace, NForm, NFormItem, NInput, NAlert,
    NGrid, NGridItem, NDivider, NRadioGroup, NRadioButton, NSwitch
} from 'naive-ui';
import { AddOutline, LogoMarkdown, BookmarkOutline } from '@vicons/ionicons5';

import { usePluginsStore } from "@/store/plugins"

import PluginCard from "./plugin-card.vue";

const show = ref(false);

watch(
    () => show.value,
    (val) => {
        val && reset();
    }
);

const form = ref();
const model = ref({
    title: "Lorem Ipsum",
    description: "",
    kind: "EMBEDDED",
    public: true,
    logo: "",
    embedded_conf: {
        frame: ""
    }
});
const rules = ref({
    title: [{ required: true, message: "Please input title" }],
});

function reset() {
    model.value = {
        title: "Lorem Ipsum",
        description: "",
        kind: "EMBEDDED",
        public: true,
        logo: "",
        embedded_conf: {
            frame: ""
        }
    };
}

const store = usePluginsStore()
const error = ref(false);
const bar = useLoadingBar();
function handleSubmit() {
    error.value = false
    form.value.validate(async (errors) => {
        if (errors) {
            return;
        }
        bar.start()

        try {
            await store.create(model.value, bar);
            bar.finish()
            show.value = false
            store.fetchPlugins()
        } catch (err) {
            bar.error()
            console.error(err)
            error.value = `${err.response.status}: ${(err.response.data ?? { message: "Unexpected Error" }).message
                }`;
        }
    });
}
</script>