<template>
    <render />
</template>

<script setup>
import { NTooltip, NTable, NCode, NIcon, NSpace, useThemeVars } from 'naive-ui';
import { defineAsyncComponent, h } from 'vue';

const IdCardOutline = defineAsyncComponent(() => import('@vicons/ionicons5/IdCardOutline'))
const TelescopeOutline = defineAsyncComponent(() => import('@vicons/ionicons5/TelescopeOutline'))

const theme = useThemeVars()
console.dir(theme.value)

const props = defineProps({
    connection: {
        type: Object, required: true
    }
})

function tooltip_row({ icon, label, value }) {
    return h('tr', [
        h('th', label),
        h('th',
            h(NIcon, { component: icon, size: 20, style: { paddingTop: '10px'  }})
        ),
        h('td', value)
    ])
}

function render() {

    let status = 'offline'
    let seen = 'Never or more than 24 hours ago'

    if (props.connection) {
        status = props.connection.connected ? 'online' : 'offline'
        if (props.connection.timestamp) {
            let d = new Date(props.connection.timestamp)
            seen = d.toString().split(' (')[0]

            if (((new Date) - d) / 1000 >= 3600) {
                status = 'onlineish'
            }
        }
    }

    let tooltip = [
        tooltip_row({
            icon: TelescopeOutline, label: 'Last Seen', value: seen
        })
    ]

    if (props.connection.clientId) {
        tooltip.push(tooltip_row({
            icon: IdCardOutline, label: 'Client ID', value: h(NCode, props.connection.clientId)
        }))
    }

    return h(NTooltip, {
        trigger: "hover",
        placement: "bottom",
    }, {
        trigger: () => h('div', {
            class: ['triangle', status]
        }),
        default: () => h(NTable, { bordered: false, bottomBordered: false }, () => h('tbody', tooltip))
    })
}
</script>

<style>
.triangle {
    width: 0;
    height: 0;
    border-style: solid;
    position: absolute;

    border-width: 0 0 30px 30px;
    left: 0;
    top: 0;
}

.online {
    border-color: transparent transparent transparent #52c41a;
}

.onlineish {
    border-color: transparent transparent transparent #f2c97d;
}

.offline {
    border-color: transparent transparent transparent var(--n-text-color);
}
</style>