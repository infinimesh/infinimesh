<template>
    <render />
</template>

<script setup>
import { NTooltip } from 'naive-ui';
import { h } from 'vue';


const props = defineProps({
    connection: {
        type: Object, required: true
    }
})


function render() {

    let status = 'offline'
    let seen = 'Never or more than 24 hours ago'

    if (props.connection) {
        status = props.connection.connected ? 'online' : 'offline'
        if (props.connection.timestamp) {
            let d = new Date(props.connection.timestamp)
            seen = d.toString().split(' (')[0]
        }
    }

    return h(NTooltip, {
        trigger: "hover",
        placement: "bottom",
    }, {
        trigger: () => h('div', {
            class: ['triangle', status]
        }),
        default: () => `Last Seen: ${seen}`
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

.offline {
    border-color: transparent transparent transparent var(--n-text-color);
}
</style>