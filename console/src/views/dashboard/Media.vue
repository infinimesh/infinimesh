<template>
  <n-grid item-responsive y-gap="10">
    <n-grid-item span="24 500:8 600:6 1000:4">
      <n-h1 prefix="bar" align-text type="info">
        <n-text type="info"> Media </n-text>
      </n-h1>
    </n-grid-item>
    <n-grid-item span="0 600:2 700:4 1000:12 1400:14"> </n-grid-item>
    <n-grid-item span="24 300:12 500:7 600:6 700:5 1000:4 1400:3">
      <n-button strong secondary round type="info" @click="stat" v-if="selected != 'all'">
        <template #icon>
          <n-icon>
            <refresh-outline />
          </n-icon>
        </template>
        Refresh
      </n-button>
    </n-grid-item>
  </n-grid>
  <n-space v-if="selected == 'all'" align="center" justify="center" class="fullscreen">
    <n-alert title="No Namespace Selected" type="info">
      <template #icon>
        <n-icon>
          <git-network-outline />
        </n-icon>
      </template>
      Media is binded to the namespace, so you must select a namespace to view and edit it's media.
      <n-space justify="center" style="margin-top: 20px">
        <ns-selector />
      </n-space>
    </n-alert>
  </n-space>
  <template v-else>
    <n-space justify="center">
      <n-upload directory-dnd @before-upload="handleBeforeUpload">
        <n-upload-dragger>
          <div style="margin-bottom: 12px">
            <n-icon size="48" :depth="3">
              <archive-outline />
            </n-icon>
          </div>
          <n-p style="font-size: 16px">
            Limit is:
            <n-number-animation :from="0" :to="limit[0]" /> {{ limit[1] }}
          </n-p>
          <n-text style="font-size: 16px">
            Click or drag a file to this area to upload
          </n-text>
          <n-p depth="3" style="margin: 8px 0 0 0">
            Strictly prohibit from uploading sensitive information. All files will be publicly accessible.
          </n-p>
          <n-p depth="3" style="margin: 8px 0 0 0">
            Also note that all special characters will be removed.
          </n-p>
        </n-upload-dragger>
      </n-upload>
    </n-space>
    <n-table :bordered="false" :single-line="false">
      <thead>
        <tr>
          <th style="width: 55%">File</th>
          <th style="width: 15%">Size</th>
          <th style="width: 15%">Last modified</th>
          <th style="width: 15%">Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="file in files" :key="file.name">
          <td @contextmenu.prevent="e => handleCopyLinkToClipboard(file.link)">
            <n-tooltip trigger="hover">
              <template #trigger>
                <a :download="file.name" :href="file.link" target="_blank" style="color: var(--n-td-text-color)">
                  {{ file.name }}
                </a>
              </template>
              Click to see the preview or download the file or the Right Click to copy the link
            </n-tooltip>
          </td>
          <td>
            <n-tooltip trigger="hover">
              <template #trigger>
                {{ sizeConv(file.size) }}
              </template>
              {{ file.size }} bytes
            </n-tooltip>
          </td>
          <td>{{ timeConf(file.mod_time) }}</td>
          <td>
            <n-popconfirm @positive-click="e => rm(file.link)">
              <template #trigger>
                <n-button round secondary type="error">
                  <template #icon>
                    <n-icon>
                      <trash-outline />
                    </n-icon>
                  </template>
                  Delete
                </n-button>
              </template>
              <span>
                Are you sure about deleting this file?
              </span>
            </n-popconfirm>
          </td>
        </tr>
      </tbody>
    </n-table>
  </template>
</template>

<script setup>
import { ref, watch, onMounted } from "vue"

import { NH1, NP, NNumberAnimation, NText, NGrid, NGridItem, NButton, NIcon, NSpace, NAlert, NTable, NTooltip, useMessage, NPopconfirm, NUpload, NUploadDragger } from 'naive-ui';
import { RefreshOutline, GitNetworkOutline, TrashOutline, ArchiveOutline } from '@vicons/ionicons5';

import { useAppStore } from '@/store/app';
import { useNSStore } from '@/store/namespaces';
import { storeToRefs } from 'pinia';

import nsSelector from '@/components/core/ns-selector.vue';

const { selected } = storeToRefs(useNSStore())
const store = useAppStore()

const files = ref([])
const limit = ref([0, "bytes"])

const base_url = `${store.console_services.http_fs}`

const units = ['bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
function sizeConv(bytes) {
  let l = 0, n = bytes;
  while (n >= 1024 && ++l) {
    n = n / 1024;
  }
  return (n.toFixed(n < 10 && l > 0 ? 1 : 0) + ' ' + units[l]);
}
function timeConf(ts) {
  return (new Date(ts)).toUTCString()
}

function makeLink(file) {
  return `${base_url}/${selected.value}/${file.name}`
}

function stat() {
  if (selected.value == 'all') {
    return
  }
  store.http.get(base_url + '/' + selected.value).then(res => {
    files.value = res.data.files.map(el => {
      el.link = makeLink(el)
      return el
    })
    let label = sizeConv(res.data.file_limit)
    limit.value = label.split(' ')
  })
}
watch(selected, stat);
function rm(link) {
  store.http.delete(link).then(res => {
    message.success('File deleted successfuly')
    stat()
  }).catch(err => {
    message.error(err.response.data.message)
  })
}

const message = useMessage()
async function handleCopyLinkToClipboard(link) {
  try {
    await navigator.clipboard.writeText(link);
    message.success("Link copied to clipboard");
  } catch {
    message.error("Failed to copy Link to clipboard");
  }
}

onMounted(() => {
  stat()
})

function handleBeforeUpload({ file }) {
  console.log(file)

  console.log(file.name.replace(/[^a-zA-Z0-9\.]/g, ""))
  let filename = file.name.replace(/[^a-zA-Z0-9\.]/g, "")

  let data = new FormData();
  data.append('file', file.file, filename);

  message.info('Uploading file...')
  store.http({
    method: 'POST',
    url: base_url + '/' + selected.value + '/' + filename,
    data
  }).then(res => {
    message.success('File uploaded successfuly')
    stat()
  }).catch(err => {
    console.log(err.response)
    message.error(err.response.data)
  })

  return false
}
</script>