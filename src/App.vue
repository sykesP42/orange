<template>
  <div id="app">
    <router-view />
    <MessageNotify />
    <ConfirmDialog />
  </div>
</template>

<script setup>
import MessageNotify from './components/MessageNotify.vue'
import ConfirmDialog from './components/ConfirmDialog.vue'
</script>

<style>
/* 全局样式保留在 style.css 和 transitions.css 中 */
</style>