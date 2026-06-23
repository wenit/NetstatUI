<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useErrorDialog } from '../composables/useErrorDialog'

const { t } = useI18n()
const { current, dismiss } = useErrorDialog()

const title = computed(() => current.value?.title ?? '')
const body = computed(() => current.value?.body ?? '')
const visible = computed(() => current.value !== null)

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape' && current.value) dismiss()
}

onMounted(() => document.addEventListener('keydown', onKey))
onUnmounted(() => document.removeEventListener('keydown', onKey))
</script>

<template>
  <transition name="fade">
    <div v-if="visible" class="overlay" @click="dismiss">
      <div class="dialog" @click.stop>
        <div class="header">
          <span class="icon">!</span>
          <span class="title">{{ title }}</span>
        </div>
        <div class="body">{{ body }}</div>
        <div class="footer">
          <button class="ok-btn" @click="dismiss">{{ t('error.ok') }}</button>
        </div>
      </div>
    </div>
  </transition>
</template>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 4000;
  backdrop-filter: blur(3px);
}
.dialog {
  width: 400px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-strong);
  border-radius: var(--radius-md);
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.4);
  overflow: hidden;
}
.header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 18px 12px;
}
.icon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--danger);
  color: #fff;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  flex-shrink: 0;
}
.title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
}
.body {
  padding: 0 18px 16px 52px;
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}
.footer {
  display: flex;
  justify-content: flex-end;
  padding: 12px 18px 16px;
  border-top: 1px solid var(--border);
  background: var(--bg);
}
.ok-btn {
  height: 32px;
  padding: 0 18px;
  border: none;
  background: var(--accent);
  color: #fff;
  font-size: 12px;
  font-weight: 600;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}
.ok-btn:hover {
  filter: brightness(1.1);
}
.ok-btn:active {
  filter: brightness(0.95);
}
.fade-enter-active, .fade-leave-active { transition: opacity var(--transition); }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.fade-enter-active .dialog, .fade-leave-active .dialog {
  transition: transform var(--transition);
}
.fade-enter-from .dialog, .fade-leave-to .dialog {
  transform: scale(0.95);
}
</style>