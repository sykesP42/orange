<template>
  <div class="workflow-page">
    <div class="page-header">
      <div class="header-content">
        <h1>工作流自动化</h1>
        <p>配置状态变更自动触发后续操作</p>
      </div>
      <div class="header-buttons">
        <button class="add-btn" @click="showCreateDialog = true">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="16"/>
            <line x1="8" y1="12" x2="16" y2="12"/>
          </svg>
          创建规则
        </button>
      </div>
    </div>

    <div class="rules-section">
      <div class="rules-list" v-if="rules.length > 0">
        <div
          v-for="rule in rules"
          :key="rule.id"
          class="rule-card"
          :class="{ disabled: !rule.enabled }"
        >
          <div class="rule-header">
            <div class="rule-info">
              <span class="rule-name">{{ rule.name }}</span>
              <span class="rule-badge" :class="rule.trigger_type">{{ getTriggerLabel(rule.trigger_type) }}</span>
            </div>
            <div class="rule-actions">
              <button class="toggle-btn" @click="toggleRule(rule)" :title="rule.enabled ? '禁用' : '启用'">
                <svg v-if="rule.enabled" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                  <line x1="1" y1="1" x2="23" y2="23"/>
                </svg>
              </button>
              <button class="edit-btn" @click="editRule(rule)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
              </button>
              <button class="delete-btn" @click="deleteRule(rule)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3,6 5,6 21,6"/>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                </svg>
              </button>
            </div>
          </div>
          <div class="rule-body">
            <div class="rule-condition">
              <span class="condition-label">触发条件：</span>
              <span class="condition-value">{{ getTriggerLabel(rule.trigger_type) }} → {{ rule.trigger_value }}</span>
            </div>
            <div class="rule-action">
              <span class="action-label">执行操作：</span>
              <span class="action-value">{{ getActionLabel(rule.action_type) }}</span>
            </div>
          </div>
          <div class="rule-footer">
            <span class="rule-priority">优先级: {{ rule.priority }}</span>
            <span class="rule-time">{{ formatTime(rule.create_time) }}</span>
          </div>
        </div>
      </div>
      <div class="empty-state" v-else>
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
        </svg>
        <p>暂无工作流规则</p>
        <button class="add-btn" @click="showCreateDialog = true">创建第一个规则</button>
      </div>
    </div>

    <div class="info-section">
      <h3>💡 工作流说明</h3>
      <div class="info-grid">
        <div class="info-item">
          <strong>触发类型：</strong>
          <span>状态变更 - 当博主状态发生变化时触发</span>
        </div>
        <div class="info-item">
          <strong>操作类型：</strong>
          <span>发送通知 / 添加标签 / 变更状态 / Webhook</span>
        </div>
        <div class="info-item">
          <strong>优先级：</strong>
          <span>数字越大优先级越高，高优先级规则先执行</span>
        </div>
        <div class="info-item">
          <strong>示例：</strong>
          <span>当状态变为"洽谈中"时，自动添加"待跟进"标签</span>
        </div>
      </div>
    </div>

    <div v-if="showCreateDialog" class="dialog-overlay" @click.self="closeDialog">
      <div class="dialog">
        <div class="dialog-header">
          <h3>{{ editingRule ? '编辑规则' : '创建工作流规则' }}</h3>
          <button class="close-btn" @click="closeDialog">×</button>
        </div>
        <div class="dialog-content">
          <div class="form-group">
            <label>规则名称</label>
            <input v-model="form.name" type="text" placeholder="例如：洽谈中自动添加标签" />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>触发类型</label>
              <select v-model="form.trigger_type">
                <option value="status_change">状态变更</option>
              </select>
            </div>
            <div class="form-group">
              <label>触发值</label>
              <select v-model="form.trigger_value">
                <option value="">请选择</option>
                <option v-for="status in statusList" :key="status.id" :value="status.name">
                  {{ status.name }}
                </option>
              </select>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>操作类型</label>
              <select v-model="form.action_type">
                <option value="notification">发送通知</option>
                <option value="tag_add">添加标签</option>
                <option value="status_change">变更状态</option>
              </select>
            </div>
            <div class="form-group">
              <label>优先级</label>
              <input v-model.number="form.priority" type="number" min="1" max="100" />
            </div>
          </div>
          <div class="form-group" v-if="form.action_type === 'notification'">
            <label>通知配置</label>
            <div class="config-fields">
              <input v-model="form.actionConfig.title" type="text" placeholder="通知标题" />
              <textarea v-model="form.actionConfig.content" placeholder="通知内容，支持变量：{nickname}" rows="3"></textarea>
            </div>
          </div>
          <div class="form-group" v-if="form.action_type === 'tag_add'">
            <label>选择标签</label>
            <div class="tag-checkboxes">
              <label v-for="tag in allTagsList" :key="tag.id" class="tag-checkbox-item">
                <input type="checkbox" v-model="form.actionConfig.tag_ids" :value="tag.id" />
                <span class="tag-name" :style="{ backgroundColor: tag.color + '20', color: tag.color }">{{ tag.name }}</span>
              </label>
            </div>
          </div>
          <div class="form-group" v-if="form.action_type === 'status_change'">
            <label>变更状态为</label>
            <select v-model="form.actionConfig.new_status">
              <option value="">请选择</option>
              <option v-for="status in statusList" :key="status.id" :value="status.name">
                {{ status.name }}
              </option>
            </select>
          </div>
          <div class="form-group checkbox-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="form.enabled" />
              <span>立即启用</span>
            </label>
          </div>
        </div>
        <div class="dialog-footer">
          <button class="btn-secondary" @click="closeDialog">取消</button>
          <button class="btn-primary" @click="saveRule" :disabled="saving">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useNotification } from '../stores/notification'
import { useConfirm } from '../utils/confirm'
import { getWorkflowRulesAPI, createWorkflowRuleAPI, updateWorkflowRuleAPI, deleteWorkflowRuleAPI, toggleWorkflowRuleAPI, getTagsAPI, getBloggerStatusList } from '../api'

const notification = useNotification()
const { confirmDanger } = useConfirm()

const rules = ref([])
const allTagsList = ref([])
const statusList = ref([])
const showCreateDialog = ref(false)
const editingRule = ref(null)
const saving = ref(false)

const form = reactive({
  name: '',
  trigger_type: 'status_change',
  trigger_value: '',
  action_type: 'notification',
  actionConfig: {
    title: '',
    content: '',
    tag_ids: [],
    new_status: ''
  },
  enabled: true,
  priority: 1
})

const loadRules = async () => {
  try {
    const res = await getWorkflowRulesAPI()
    if (res.code === 200) {
      rules.value = res.data || []
    }
  } catch (error) {
    console.error('加载工作流规则失败', error)
  }
}

const loadTags = async () => {
  try {
    const res = await getTagsAPI()
    if (res.code === 200) {
      allTagsList.value = res.data || []
    }
  } catch (error) {
    console.error('加载标签失败', error)
  }
}

const loadStatusList = async () => {
  try {
    const res = await getBloggerStatusList()
    if (res.code === 200) {
      statusList.value = res.data || []
    }
  } catch (error) {
    console.error('加载状态列表失败', error)
  }
}

const getTriggerLabel = (type) => {
  const labels = {
    'status_change': '状态变更'
  }
  return labels[type] || type
}

const getActionLabel = (type) => {
  const labels = {
    'notification': '发送通知',
    'tag_add': '添加标签',
    'status_change': '变更状态',
    'email': '发送邮件',
    'webhook': 'Webhook'
  }
  return labels[type] || type
}

const formatTime = (time) => {
  if (!time) return ''
  const d = new Date(time)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const toggleRule = async (rule) => {
  try {
    const res = await toggleWorkflowRuleAPI(rule.id)
    if (res.code === 200) {
      rule.enabled = rule.enabled ? 0 : 1
      notification.success(rule.enabled ? '规则已启用' : '规则已禁用')
    }
  } catch (error) {
    notification.error('操作失败')
  }
}

const editRule = (rule) => {
  editingRule.value = rule
  form.name = rule.name
  form.trigger_type = rule.trigger_type
  form.trigger_value = rule.trigger_value
  form.action_type = rule.action_type
  form.enabled = rule.enabled === 1
  form.priority = rule.priority

  try {
    form.actionConfig = JSON.parse(rule.action_config || '{}')
  } catch {
    form.actionConfig = { title: '', content: '', tag_ids: [], new_status: '' }
  }

  showCreateDialog.value = true
}

const deleteRule = async (rule) => {
  if (!await confirmDanger(`确定要删除规则「${rule.name}」吗？`)) return

  try {
    const res = await deleteWorkflowRuleAPI(rule.id)
    if (res.code === 200) {
      notification.success('删除成功')
      loadRules()
    }
  } catch (error) {
    notification.error('删除失败')
  }
}

const closeDialog = () => {
  showCreateDialog.value = false
  editingRule.value = null
  form.name = ''
  form.trigger_type = 'status_change'
  form.trigger_value = ''
  form.action_type = 'notification'
  form.actionConfig = { title: '', content: '', tag_ids: [], new_status: '' }
  form.enabled = true
  form.priority = 1
}

const saveRule = async () => {
  if (!form.name.trim()) {
    notification.warning('请输入规则名称')
    return
  }
  if (!form.trigger_value) {
    notification.warning('请选择触发值')
    return
  }

  saving.value = true
  try {
    const config = { ...form.actionConfig }
    if (form.action_type === 'tag_add') {
      config.tag_ids = form.actionConfig.tag_ids || []
    }

    const data = {
      name: form.name,
      trigger_type: form.trigger_type,
      trigger_value: form.trigger_value,
      action_type: form.action_type,
      action_config: JSON.stringify(config),
      enabled: form.enabled,
      priority: form.priority
    }

    let res
    if (editingRule.value) {
      res = await updateWorkflowRuleAPI(editingRule.value.id, data)
    } else {
      res = await createWorkflowRuleAPI(data)
    }

    if (res.code === 200) {
      notification.success(editingRule.value ? '更新成功' : '创建成功')
      closeDialog()
      loadRules()
    } else {
      notification.error(res.message || '保存失败')
    }
  } catch (error) {
    notification.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadRules()
  loadTags()
  loadStatusList()
})
</script>

<style scoped>
.workflow-page {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.page-header p {
  font-size: 14px;
  color: var(--text-muted);
  margin: 8px 0 0;
}

.header-buttons {
  display: flex;
  gap: 12px;
}

.add-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 10px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.add-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.4);
}

.add-btn svg {
  width: 18px;
  height: 18px;
}

.rules-section {
  margin-bottom: 32px;
}

.rules-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 20px;
}

.rule-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 20px;
  transition: all 0.2s;
}

.rule-card:hover {
  border-color: var(--primary);
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.1);
}

.rule-card.disabled {
  opacity: 0.6;
}

.rule-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.rule-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.rule-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.rule-badge {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.rule-badge.status_change {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.rule-actions {
  display: flex;
  gap: 8px;
}

.rule-actions button {
  width: 32px;
  height: 32px;
  padding: 0;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.rule-actions button svg {
  width: 16px;
  height: 16px;
  color: var(--text-secondary);
}

.toggle-btn:hover {
  background: var(--bg-hover);
}

.edit-btn:hover {
  border-color: #667eea;
}

.edit-btn:hover svg {
  color: #667eea;
}

.delete-btn:hover {
  border-color: var(--danger);
  background: rgba(239, 68, 68, 0.05);
}

.delete-btn:hover svg {
  color: var(--danger);
}

.rule-body {
  margin-bottom: 16px;
}

.rule-condition,
.rule-action {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-bottom: 8px;
  font-size: 14px;
}

.condition-label,
.action-label {
  color: var(--text-muted);
  min-width: 70px;
}

.condition-value,
.action-value {
  color: var(--text-primary);
}

.rule-footer {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--text-muted);
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: var(--bg-card);
  border-radius: 12px;
  border: 1px dashed var(--border-color);
}

.empty-state svg {
  width: 48px;
  height: 48px;
  color: var(--text-muted);
  margin-bottom: 16px;
}

.empty-state p {
  color: var(--text-muted);
  margin-bottom: 20px;
}

.info-section {
  background: var(--bg-card);
  border-radius: 12px;
  padding: 24px;
}

.info-section h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.info-item {
  font-size: 14px;
  color: var(--text-secondary);
}

.info-item strong {
  color: var(--text-primary);
}

.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
  background: var(--bg-card);
  border-radius: 16px;
  width: 90%;
  max-width: 560px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
}

.dialog-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.close-btn {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: transparent;
  border: none;
  font-size: 20px;
  color: var(--text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.close-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.dialog-content {
  padding: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 12px 14px;
  background: var(--bg-dark);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  color: var(--text-primary);
  transition: border-color 0.2s;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #667eea;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.config-fields {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.tag-checkboxes {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  max-height: 150px;
  overflow-y: auto;
}

.tag-checkbox-item {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.tag-checkbox-item input {
  width: 16px;
  height: 16px;
  cursor: pointer;
  accent-color: #667eea;
}

.tag-checkbox-item .tag-name {
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  border: 1px solid;
}

.checkbox-group {
  margin-top: 24px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  color: var(--text-primary);
}

.checkbox-label input {
  width: 16px;
  height: 16px;
  cursor: pointer;
  accent-color: #667eea;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid var(--border-color);
}

.btn-secondary {
  padding: 10px 20px;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-secondary);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.btn-primary {
  padding: 10px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .rules-list {
    grid-template-columns: 1fr;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>