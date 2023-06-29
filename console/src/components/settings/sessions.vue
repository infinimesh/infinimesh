<template>
    <n-space justify="space-between">
        <n-h2 prefix="bar">
            <n-text>
                Sessions
            </n-text>
        </n-h2>

        <n-button round type="info" ghost @click="load">
            <template #icon>
                <n-icon :component="RefreshOutline" />
            </template>
            Refresh
        </n-button>
    </n-space>
    <n-table>
        <thead>
            <tr>
                <th colspan="2">ID</th>
                <th>Client</th>
                <th>Signed In</th>
                <th>Last Seen</th>
                <th>Expires</th>
                <th></th>
            </tr>
        </thead>
        <tbody>
            <template v-if="loading">
                <tr v-for="i in 5">
                    <th colspan="2"><n-skeleton height="2rem" :sharp="false" /></th>
                    <th><n-skeleton height="2rem" :sharp="false" /></th>
                    <th><n-skeleton height="2rem" :sharp="false" /></th>
                    <th><n-skeleton height="2rem" :sharp="false" /></th>
                    <th><n-skeleton height="2rem" :sharp="false" /></th>
                </tr>
            </template>
            <template v-else>
                <tr v-for="session in  sessions ">
                    <template v-if="session.current">
                        <th>
                            <n-tooltip trigger="hover">
                                <template #trigger>
                                    <n-code :style="{ color: theme.successColor }">{{ session.id }}</n-code>
                                </template>
                                This is your current session.
                            </n-tooltip>
                        </th>
                        <th>
                            <n-icon :component="ArrowDownCircleOutline" size="24" style="padding-top: 8px;"
                                :color="theme.successColor" />
                        </th>
                    </template>
                    <th v-else colspan="2" align="right"><n-code>{{ session.id }}</n-code></th>
                    <th>
                        <client :client="session.client ? session.client : 'Unknown'" />
                    </th>
                    <td>
                        <relative_time :timestamp="session.created" />
                    </td>

                    <td v-if="activity[session.id]">
                        <relative_time :timestamp="activity[session.id]" />
                    </td>
                    <td v-else>
                        Never
                    </td>

                    <td v-if="session.expires">
                        <relative_time :timestamp="session.expires" />
                    </td>
                    <td v-else>
                        <n-text type="warning">
                            <n-icon :component="WarningOutline" />
                            Never
                        </n-text>
                    </td>

                    <td>
                        <n-tooltip trigger="hover" v-if="!session.current">
                            <template #trigger>
                                <n-button ghost type="warning"
                                    @click="async () => { await store.revoke(session.id); load() }">
                                    <template #icon>
                                        <n-icon :component="ExitOutline" />
                                    </template>
                                </n-button>
                            </template>

                            Revoke Session
                        </n-tooltip>
                    </td>
                </tr>
            </template>
        </tbody>
    </n-table>
</template>

<script setup>
import { ref, onMounted, defineAsyncComponent, h } from 'vue';
import { useSessionsStore } from '@/store/sessions';

import {
    NTable, NSkeleton, NCode,
    NIcon, useThemeVars, NTooltip,
    NText, NDivider, NButton,
    NSpace, NH2,
} from 'naive-ui'

const ArrowDownCircleOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ArrowDownCircleOutline'))
const ExitOutline = defineAsyncComponent(() => import('@vicons/ionicons5/ExitOutline'))
const RefreshOutline = defineAsyncComponent(() => import('@vicons/ionicons5/RefreshOutline'))
const WarningOutline = defineAsyncComponent(() => import('@vicons/ionicons5/WarningOutline'))

const theme = useThemeVars()

const store = useSessionsStore()
const loading = ref(false)

const sessions = ref([])
const activity = ref({})

async function load() {
    loading.value = true
    await Promise.all([
        new Promise((resolve) => {
            store.get().then((sess) => {
                for (let i = 0; i < sess.sessions.length; i++) {
                    if (sess.sessions[i].current) {
                        let s = sess.sessions[i]
                        sess.sessions.splice(i, 1)
                        sess.sessions.unshift(s)
                    }
                }

                sessions.value = sess.sessions
                resolve()
            })
        }),
        new Promise((resolve) => {
            store.activity().then((act) => {
                activity.value = act.lastSeen
                resolve()
            })
        })
    ])
    loading.value = false
}

onMounted(load)

const now = ref(new Date())
setInterval(() => now.value = new Date(), 1000)

function relative_time({ timestamp }) {
    timestamp = new Date(timestamp)
    let label = timestamp < now.value.getTime() ? relative_time_past(timestamp) : relative_time_future(timestamp)

    return h(NTooltip, { trigger: 'hover' }, {
        default: () => timestamp.toLocaleDateString(),
        trigger: () => label
    })
}

function relative_time_past(timestamp) {
    const timeDifference = (now.value.getTime() - timestamp.getTime()) / 1000;
    const minutesDifference = Math.floor(timeDifference / 60);

    if (minutesDifference >= 4320) {
        return timestamp.toLocaleDateString();
    } else if (minutesDifference >= 1440) {
        const daysDifference = Math.floor(minutesDifference / 1440);
        return `${daysDifference} days ago`;
    } else if (minutesDifference >= 60) {
        const hoursDifference = Math.floor(minutesDifference / 60);
        return `${hoursDifference} hours ago`;
    } else if (minutesDifference > 0) {
        return `${minutesDifference} minutes ago`;
    } else {
        return 'just now';
    }
}

function relative_time_future(timestamp) {
    const timeDifference = (timestamp.getTime() - now.value.getTime()) / 1000;
    const minutesDifference = Math.floor(timeDifference / 60);

    if (minutesDifference >= 4320) {
        return timestamp.toLocaleDateString();
    } else if (minutesDifference >= 1440) {
        const daysDifference = Math.floor(minutesDifference / 1440);
        return `in ${daysDifference} days`;
    } else if (minutesDifference >= 60) {
        const hoursDifference = Math.floor(minutesDifference / 60);
        return `in ${hoursDifference} hours`;
    } else if (minutesDifference > 0) {
        return `in ${minutesDifference} minutes`;
    } else {
        return 'in few seconds';
    }
}

const icons = {
    'chrome': defineAsyncComponent(() => import('@vicons/ionicons5/LogoChrome')),
    'firefox': defineAsyncComponent(() => import('@vicons/ionicons5/LogoFirefox')),
    'linux': defineAsyncComponent(() => import('@vicons/ionicons5/LogoTux')),
    'mac': defineAsyncComponent(() => import('@vicons/ionicons5/LogoApple')),
    'windows': defineAsyncComponent(() => import('@vicons/ionicons5/LogoWindows')),
    'android': defineAsyncComponent(() => import('@vicons/ionicons5/LogoAndroid')),
    'cli': defineAsyncComponent(() => import('@vicons/ionicons5/TerminalOutline')),
    'unknown': defineAsyncComponent(() => import('@vicons/ionicons5/HelpOutline')),
    'console': defineAsyncComponent(() => import('@vicons/ionicons5/DesktopOutline'))
}
const interleave = (arr, thing) => [].concat(...arr.map(n => [n, thing])).slice(0, -1)

function client({ client }) {

    let parts = client.split('|')

    parts = parts.map((part, i) => {
        part = part.trim()

        let icon = icons[part.toLowerCase()]
        if (icon)
            part = [h(NIcon, { component: icon, size: '1rem', style: 'margin-right: 4px;' }), part]

        return h(NText, { depth: i > 2 ? 3 : i + 1, style: { fontSize: '1.2rem' } }, () => part)
    })

    return interleave(parts, h(NDivider, { vertical: true, style: 'margin: 0 4px;' }))
}

</script>