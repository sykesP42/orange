<template>
  <div class="forums-page">
    <div class="page-header">
      <div class="header-left">
        <button class="back-btn" @click="$router.back()" v-if="showBack">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="19" y1="12" x2="5" y2="12"/>
            <polyline points="12 19 5 12 12 5"/>
          </svg>
        </button>
        <div class="header-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
          </svg>
        </div>
        <div class="header-text">
          <h2>公共论坛</h2>
          <p class="page-desc">与所有用户一起交流讨论</p>
        </div>
      </div>
      <div class="header-right">
        <div class="header-stats">
          <div class="stat-item">
            <span class="stat-number">{{ totalPosts }}</span>
            <span class="stat-label">帖子</span>
          </div>
        </div>
        <button class="btn btn-primary create-post-header-btn" @click="showCreatePost" v-if="canPost">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          发布话题
        </button>
      </div>
    </div>

    <div class="forum-toolbar">
      <div class="search-box">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/>
          <line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <input type="text" v-model="searchKeyword" placeholder="搜索帖子..." @keyup.enter="searchPosts" />
        <button v-if="searchKeyword" class="clear-search" @click="clearSearch">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
      <div class="view-tabs">
        <button class="view-tab" :class="{ active: viewMode === 'all' }" @click="switchView('all')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
          </svg>
          全部
        </button>
        <button class="view-tab" :class="{ active: viewMode === 'hot' }" @click="switchView('hot')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2c-5.33 4.55-8 8.48-8 11.8 0 4.98 3.8 8.2 8 8.2s8-3.22 8-8.2c0-3.32-2.67-7.25-8-11.8z"/>
          </svg>
          热门
        </button>
        <button class="view-tab" :class="{ active: viewMode === 'collected' }" @click="switchView('collected')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"/>
          </svg>
          收藏
        </button>
      </div>
    </div>

    <div class="category-tabs">
      <button
        v-for="cat in categories"
        :key="cat"
        class="category-tab"
        :class="{ active: selectedCategory === cat }"
        @click="selectCategory(cat)"
      >
        {{ cat }}
      </button>
    </div>

    <div class="posts-container" v-if="posts.length > 0">
      <div v-for="post in posts" :key="post.id" class="post-card" :class="{ pinned: post.is_pinned, featured: post.is_featured }" @click="viewPost(post)">
        <div class="post-badges" v-if="post.is_pinned || post.is_featured">
          <span v-if="post.is_pinned" class="badge pinned">
            <svg viewBox="0 0 24 24" fill="currentColor" width="12" height="12">
              <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/>
            </svg>
            置顶
          </span>
          <span v-if="post.is_featured" class="badge featured">
            <svg viewBox="0 0 24 24" fill="currentColor" width="12" height="12">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
            </svg>
            精华
          </span>
        </div>
        <div class="post-header">
          <div class="post-author">
            <div class="author-avatar" :style="{ backgroundColor: post.author_color || 'var(--primary)' }">
              <img v-if="post.author_avatar" :src="post.author_avatar" :alt="post.author_name" v-avatar />
              <span v-else>{{ post.author_name?.charAt(0) || 'U' }}</span>
            </div>
            <div class="author-info">
              <span class="author-name">{{ post.author_name }}</span>
              <span class="post-time">{{ formatTime(post.create_time) }}</span>
            </div>
          </div>
          <span class="post-category" :class="post.category">{{ post.category }}</span>
        </div>
        <h3 class="post-title">{{ post.title }}</h3>
        <p class="post-content">{{ truncateContent(post.content) }}</p>
        <div class="post-footer">
          <div class="post-stats">
            <span class="stat-item">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
              {{ post.view_count || 0 }}
            </span>
            <span class="stat-item">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
              </svg>
              {{ post.like_count || 0 }}
            </span>
            <span class="stat-item">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
              </svg>
              {{ post.comment_count || 0 }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="empty-state">
      <div class="empty-icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
        </svg>
      </div>
      <p>暂无帖子</p>
      <span class="empty-hint">成为第一个发表帖子的人吧</span>
      <button class="btn btn-primary" @click="showCreatePost" v-if="canPost">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"/>
          <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        发表帖子
      </button>
    </div>

    <div class="load-more" v-if="hasMore" @click="loadMore">
      <button class="btn btn-secondary">加载更多</button>
    </div>

    <div v-if="selectedPost" class="modal-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="closePost">
      <div class="modal" style="background:var(--bg-card);max-width:640px;">
        <div class="detail-header">
          <div class="detail-title-row">
            <div class="post-badges" v-if="selectedPost.is_pinned || selectedPost.is_featured">
              <span v-if="selectedPost.is_pinned" class="badge pinned">置顶</span>
              <span v-if="selectedPost.is_featured" class="badge featured">精华</span>
            </div>
            <h3>{{ selectedPost.title }}</h3>
          </div>
          <button class="close-btn" @click="closePost">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="detail-body">
          <div class="post-meta">
            <div class="author-avatar" :style="{ backgroundColor: selectedPost.author_color || 'var(--purple)' }">
              <img v-if="selectedPost.author_avatar" :src="selectedPost.author_avatar" :alt="selectedPost.author_name" v-avatar />
              <span v-else>{{ selectedPost.author_name?.charAt(0) || 'U' }}</span>
            </div>
            <div class="meta-info">
              <span class="author-name">{{ selectedPost.author_name }}</span>
              <span class="post-time">{{ formatDate(selectedPost.create_time) }}</span>
            </div>
            <span class="post-category" :class="selectedPost.category">{{ selectedPost.category }}</span>
          </div>
          <div class="post-content-detail">
            <p>{{ selectedPost.content }}</p>
          </div>
          <div class="post-stats-bar">
            <span class="stat-item">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
              {{ selectedPost.view_count || 0 }} 浏览
            </span>
            <span class="stat-item">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
              </svg>
              {{ selectedPost.like_count || 0 }} 赞
            </span>
            <span class="stat-item">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
              </svg>
              {{ selectedPost.comment_count || 0 }} 评论
            </span>
          </div>
          <div class="post-actions-detail">
            <button class="action-btn" :class="{ active: hasLiked }" @click="likePost">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
              </svg>
              {{ hasLiked ? '已赞' : '点赞' }}
            </button>
            <button class="action-btn" :class="{ active: hasCollected }" @click="collectPost">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"/>
              </svg>
              {{ hasCollected ? '已收藏' : '收藏' }}
            </button>
            <button class="action-btn" @click="showCommentInput = true">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
              </svg>
              评论
            </button>
            <button v-if="isAdmin" class="action-btn" :class="{ active: selectedPost.is_pinned }" @click="pinPost">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/>
              </svg>
              {{ selectedPost.is_pinned ? '取消置顶' : '置顶' }}
            </button>
            <button v-if="isAdmin" class="action-btn" :class="{ active: selectedPost.is_featured }" @click="featurePost">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
              </svg>
              {{ selectedPost.is_featured ? '取消精华' : '精华' }}
            </button>
            <button v-if="canDelete" class="action-btn danger" @click="deletePost">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6"/>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
              </svg>
              删除
            </button>
          </div>

          <div v-if="showCommentInput" class="comment-input-section">
            <textarea v-model="commentContent" placeholder="写下你的评论..." rows="3"></textarea>
            <div class="comment-actions">
              <button class="btn btn-secondary" @click="cancelComment">取消</button>
              <button class="btn btn-primary" @click="submitComment">发表评论</button>
            </div>
          </div>

          <div class="comments-section">
            <h3>评论 ({{ comments.length }})</h3>
            <div v-if="comments.length > 0" class="comments-list">
              <div v-for="comment in comments" :key="comment.id" class="comment-item">
                <div class="comment-avatar" :style="{ backgroundColor: 'var(--purple)' }">
                  <img v-if="comment.author_avatar" :src="comment.author_avatar" :alt="comment.author_name" v-avatar />
                  <span v-else>{{ comment.author_name?.charAt(0) || 'U' }}</span>
                </div>
                <div class="comment-body">
                  <div class="comment-header">
                    <span class="comment-author">{{ comment.author_name }}</span>
                    <span class="comment-time">{{ formatTime(comment.create_time) }}</span>
                  </div>
                  <p class="comment-text">{{ comment.content }}</p>
                </div>
              </div>
            </div>
            <div v-else class="no-comments">
              暂无评论，来发表第一条评论吧
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="createPostVisible" class="modal-overlay" style="background:rgba(0,0,0,0.6)!important;backdrop-filter:blur(4px);" @click.self="closeCreatePost">
      <div class="modal" style="background:var(--bg-card);max-width:640px;">
        <div class="detail-header">
          <h3>发布帖子</h3>
          <button class="close-btn" @click="closeCreatePost">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="detail-body">
          <div class="form-group">
            <label>标题</label>
            <input v-model="newPost.title" placeholder="输入帖子标题" />
          </div>
          <div class="form-group">
            <label>分类</label>
            <select v-model="newPost.category">
              <option v-for="cat in categories.slice(1)" :key="cat" :value="cat">{{ cat }}</option>
            </select>
          </div>
          <div class="form-group">
            <label>内容</label>
            <textarea v-model="newPost.content" placeholder="输入帖子内容" rows="6"></textarea>
          </div>
        </div>
        <div class="detail-footer">
          <button class="btn btn-secondary" @click="closeCreatePost">取消</button>
          <button class="btn btn-primary" @click="submitPost">发布</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useNotification } from '../stores/notification'
import { useConfirm } from '../utils/confirm'
import { getForumPostsAPI, getForumPostsHotAPI, getForumPostsCollectedAPI, searchForumPostsAPI, getForumPostDetailAPI, getForumPostLikeStatusAPI, getForumPostCollectStatusAPI, createForumPostAPI, createForumCommentAPI, likeForumPostAPI, collectForumPostAPI, deleteForumPostAPI, pinForumPostAPI, featureForumPostAPI } from '../api'

const router = useRouter()
const userStore = useUserStore()
const notification = useNotification()
const { confirmDanger } = useConfirm()

const posts = ref([])
const page = ref(1)
const pageSize = 20
const totalPosts = ref(0)
const hasMore = computed(() => posts.value.length < totalPosts.value)
const selectedCategory = ref('全部')
const categories = ['全部', '讨论', '分享', '求助', '公告']
const createPostVisible = ref(false)
const showBack = ref(false)
const selectedPost = ref(null)
const comments = ref([])
const showCommentInput = ref(false)
const commentContent = ref('')
const hasLiked = ref(false)
const hasCollected = ref(false)
const searchKeyword = ref('')
const viewMode = ref('all')

const newPost = ref({
  title: '',
  content: '',
  category: '讨论'
})

const canPost = computed(() => !!userStore.id)
const isAdmin = computed(() => userStore.role === 'admin')
const canDelete = computed(() => {
  if (!selectedPost.value || !userStore.id) return false
  return selectedPost.value.author_id === userStore.id || userStore.role === 'admin'
})

const loadPosts = async () => {
  try {
    let res
    if (viewMode.value === 'hot') {
      res = await getForumPostsHotAPI({ limit: 50 })
    } else if (viewMode.value === 'collected') {
      res = await getForumPostsCollectedAPI({ page: page.value, pageSize })
    } else {
      res = await getForumPostsAPI({ page: page.value, pageSize, category: selectedCategory.value })
    }
    if (res.code === 200) {
      if (page.value === 1) {
        posts.value = res.data || []
      } else {
        posts.value = [...posts.value, ...(res.data || [])]
      }
      totalPosts.value = res.total || posts.value.length
    }
  } catch (error) {
    notification.error('加载帖子失败')
  }
}

const loadMore = () => {
  page.value++
  loadPosts()
}

const selectCategory = (cat) => {
  selectedCategory.value = cat
  page.value = 1
  loadPosts()
}

const switchView = (mode) => {
  viewMode.value = mode
  page.value = 1
  loadPosts()
}

const searchPosts = async () => {
  if (!searchKeyword.value.trim()) {
    loadPosts()
    return
  }
  try {
    const res = await searchForumPostsAPI({ keyword: searchKeyword.value, page: 1, pageSize })
    if (res.code === 200) {
      posts.value = res.data || []
      totalPosts.value = res.total || 0
      viewMode.value = 'all'
    }
  } catch (error) {
    notification.error('搜索失败')
  }
}

const clearSearch = () => {
  searchKeyword.value = ''
  loadPosts()
}

const viewPost = async (post) => {
  try {
    const res = await getForumPostDetailAPI(post.id)
    if (res.code === 200) {
      selectedPost.value = res.data
      comments.value = res.data.comments || []
      hasLiked.value = false
      hasCollected.value = false
      
      const likeRes = await getForumPostLikeStatusAPI(post.id)
      if (likeRes.code === 200) {
        hasLiked.value = likeRes.data?.liked || false
      }
      
      const collectRes = await getForumPostCollectStatusAPI(post.id)
      if (collectRes.code === 200) {
        hasCollected.value = collectRes.data?.collected || false
      }
    }
  } catch (error) {
    notification.error('加载帖子详情失败')
  }
}

const closePost = () => {
  selectedPost.value = null
  comments.value = []
  showCommentInput.value = false
  commentContent.value = ''
  hasLiked.value = false
  hasCollected.value = false
}

const likePost = async () => {
  if (!selectedPost.value) return
  try {
    const res = await likeForumPostAPI(selectedPost.value.id)
    if (res.code === 200) {
      selectedPost.value.like_count = res.like_count
      hasLiked.value = res.liked
    }
  } catch (error) {
    notification.error('点赞失败')
  }
}

const collectPost = async () => {
  if (!selectedPost.value) return
  try {
    const res = await collectForumPostAPI(selectedPost.value.id)
    if (res.code === 200) {
      hasCollected.value = res.collected
      notification.success(res.message)
    }
  } catch (error) {
    notification.error('收藏失败')
  }
}

const pinPost = async () => {
  if (!selectedPost.value) return
  try {
    const res = await pinForumPostAPI(selectedPost.value.id)
    if (res.code === 200) {
      selectedPost.value.is_pinned = res.is_pinned
      notification.success(res.message)
      loadPosts()
    }
  } catch (error) {
    notification.error('操作失败')
  }
}

const featurePost = async () => {
  if (!selectedPost.value) return
  try {
    const res = await featureForumPostAPI(selectedPost.value.id)
    if (res.code === 200) {
      selectedPost.value.is_featured = res.is_featured
      notification.success(res.message)
      loadPosts()
    }
  } catch (error) {
    notification.error('操作失败')
  }
}

const cancelComment = () => {
  showCommentInput.value = false
  commentContent.value = ''
}

const submitComment = async () => {
  if (!commentContent.value.trim()) {
    notification.warning('请输入评论内容')
    return
  }
  try {
    const res = await createForumCommentAPI(selectedPost.value.id, { content: commentContent.value })
    if (res.code === 200) {
      comments.value.push(res.data)
      selectedPost.value.comment_count = (selectedPost.value.comment_count || 0) + 1
      showCommentInput.value = false
      commentContent.value = ''
      notification.success('评论成功')
    } else {
      notification.error(res.message || '评论失败')
    }
  } catch (error) {
    notification.error('评论失败')
  }
}

const deletePost = async () => {
  const confirmed = await confirmDanger('确定要删除这篇帖子吗？此操作不可恢复。')
  if (!confirmed) return
  try {
    const res = await deleteForumPostAPI(selectedPost.value.id)
    if (res.code === 200) {
      notification.success('删除成功')
      closePost()
      loadPosts()
    } else {
      notification.error(res.message || '删除失败')
    }
  } catch (error) {
    notification.error('删除失败')
  }
}

const showCreatePost = () => {
  newPost.value = { title: '', content: '', category: '讨论' }
  createPostVisible.value = true
}

const closeCreatePost = () => {
  createPostVisible.value = false
}

const submitPost = async () => {
  if (!newPost.value.title.trim()) {
    notification.warning('请输入帖子标题')
    return
  }
  if (!newPost.value.content.trim()) {
    notification.warning('请输入帖子内容')
    return
  }

  try {
    const res = await createForumPostAPI(newPost.value)
    if (res.code === 200) {
      notification.success('发布成功')
      closeCreatePost()
      loadPosts()
    } else {
      notification.error(res.message || '发布失败')
    }
  } catch (error) {
    notification.error(error?.response?.data?.message || error?.message || '发布失败')
  }
}

const formatTime = (time) => {
  if (!time) return ''
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  return `${Math.floor(diff / 86400000)}天前`
}

const formatDate = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

const truncateContent = (content) => {
  if (!content) return ''
  return content.length > 120 ? content.substring(0, 120) + '...' : content
}

onMounted(() => {
  if (window.history.length > 1) {
    showBack.value = true
  }
  loadPosts()
})
</script>

<style scoped>
.forums-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 24px;
  padding-bottom: 80px;
}

.detail-header {
  padding: 20px 24px 16px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
}

.detail-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.detail-body {
  padding: 20px 24px;
  flex: 1;
  overflow-y: auto;
  background: var(--bg-card);
}

.detail-footer {
  padding: 16px 24px;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  flex-shrink: 0;
  background: var(--bg-card);
}


.page-header {
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  flex-wrap: wrap;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.back-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: 1px solid var(--border-color);
  background: var(--bg-card);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.back-btn:hover {
  background: var(--bg-hover);
}

.back-btn svg {
  width: 18px;
  height: 18px;
  color: var(--text-secondary);
}

.header-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: var(--bg-brand-gradient);
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-icon svg {
  width: 24px;
  height: 24px;
  color: var(--color-on-brand);
}

.header-text h2 {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 4px 0;
}

.page-desc {
  color: var(--text-muted);
  font-size: 14px;
  margin: 0;
}

.header-stats {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-stats .stat-item {
  background: var(--bg-card);
  padding: 8px 16px;
  border-radius: 20px;
  text-align: center;
  border: 1px solid var(--border-color);
}

.header-stats .stat-number {
  font-size: 20px;
  font-weight: 700;
  color: var(--primary);
  display: block;
}

.header-stats .stat-label {
  font-size: 12px;
  color: var(--text-muted);
}

.create-post-header-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  font-size: 14px;
  font-weight: 500;
  border-radius: 10px;
  white-space: nowrap;
}

.create-post-header-btn svg {
  width: 18px;
  height: 18px;
}

.forum-toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.search-box {
  flex: 1;
  min-width: 200px;
  position: relative;
  display: flex;
  align-items: center;
}

.search-box svg {
  position: absolute;
  left: 12px;
  width: 18px;
  height: 18px;
  color: var(--text-muted);
}

.search-box input {
  width: 100%;
  padding: 10px 36px;
  border-radius: 10px;
  border: 1px solid var(--border-color);
  background: var(--bg-card);
  color: var(--text-primary);
  font-size: 14px;
  transition: all 0.2s;
}

.search-box input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(255, 107, 53, 0.1);
}

.search-box input::placeholder {
  color: var(--text-muted);
}

.clear-search {
  position: absolute;
  right: 8px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: none;
  background: var(--bg-hover);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.clear-search svg {
  position: static;
  width: 14px;
  height: 14px;
  color: var(--text-muted);
}

.view-tabs {
  display: flex;
  gap: 4px;
  background: var(--bg-card);
  padding: 4px;
  border-radius: 10px;
  border: 1px solid var(--border-color);
}

.view-tab {
  padding: 8px 16px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  transition: all 0.2s;
}

.view-tab svg {
  width: 16px;
  height: 16px;
}

.view-tab:hover {
  background: var(--bg-hover);
}

.view-tab.active {
  background: var(--primary);
  color: var(--color-on-brand);
}

.category-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
  overflow-x: auto;
  padding-bottom: 8px;
}

.category-tab {
  padding: 8px 16px;
  border-radius: 20px;
  border: 1px solid var(--border-color);
  background: var(--bg-card);
  color: var(--text-secondary);
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s;
  font-size: 14px;
}

.category-tab:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.category-tab.active {
  background: var(--bg-brand-gradient);
  color: var(--color-on-brand);
  border-color: var(--primary);
}

.posts-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.post-card {
  background: var(--bg-card);
  border-radius: 16px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid var(--border-color);
  position: relative;
}

.post-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--primary-light);
}

.post-card.pinned {
  border-left: 4px solid var(--primary);
}

.post-card.featured {
  border-left: 4px solid #f59e0b;
}

.post-badges {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.badge.pinned {
  background: rgba(255, 107, 53, 0.1);
  color: var(--primary);
}

.badge.featured {
  background: rgba(245, 158, 11, 0.1);
  color: var(--warning);
}

.post-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.post-author {
  display: flex;
  align-items: center;
  gap: 12px;
}

.author-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: var(--bg-brand-gradient);
  color: var(--color-on-brand);
  font-weight: 600;
  font-size: 16px;
}

.author-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.author-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.author-name {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 14px;
}

.post-time {
  font-size: 12px;
  color: var(--text-muted);
}

.post-category {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  background: var(--bg-hover);
  color: var(--text-secondary);
}

.post-category.讨论 {
  background: rgba(99, 102, 241, 0.1);
  color: var(--purple);
}

.post-category.求助 {
  background: rgba(239, 68, 68, 0.1);
  color: var(--danger);
}

.post-category.分享 {
  background: rgba(16, 185, 129, 0.1);
  color: var(--success);
}

.post-category.公告 {
  background: rgba(245, 158, 11, 0.1);
  color: var(--warning);
}

.post-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  line-height: 1.4;
}

.post-content {
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.6;
  margin: 0 0 16px 0;
}

.post-footer {
  display: flex;
  justify-content: flex-end;
}

.post-stats {
  display: flex;
  gap: 16px;
}

.post-stats .stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--text-muted);
  font-size: 13px;
}

.post-stats .stat-item svg {
  width: 16px;
  height: 16px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: var(--bg-card);
  border-radius: 16px;
  border: 1px solid var(--border-color);
}

.empty-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  border-radius: 50%;
  background: var(--bg-hover);
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-icon svg {
  width: 40px;
  height: 40px;
  color: var(--text-muted);
}

.empty-state p {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px 0;
}

.empty-hint {
  color: var(--text-muted);
  font-size: 14px;
  display: block;
  margin-bottom: 20px;
}

.load-more {
  text-align: center;
  padding: 20px;
}

@media (max-width: 768px) {
  .forums-page {
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .header-right {
    justify-content: space-between;
  }

  .forum-toolbar {
    flex-direction: column;
    gap: 8px;
  }

  .search-box {
    min-width: 100%;
  }

  .view-tabs {
    justify-content: center;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    scrollbar-width: none;
  }

  .view-tabs::-webkit-scrollbar {
    display: none;
  }

  .post-title {
    font-size: 16px;
  }

  .create-post-header-btn {
    padding: 8px 16px;
    font-size: 13px;
  }

  .post-header {
    flex-direction: column;
    gap: 8px;
  }

  .post-footer {
    justify-content: flex-start;
  }
}

@media (max-width: 640px) {
  .forums-page {
    padding: 12px;
  }

  .page-header h1 {
    font-size: 22px;
  }

  .post-title {
    font-size: 15px;
  }
}

@media (max-width: 480px) {
  .forums-page {
    padding: 10px;
  }

  .page-header h1 {
    font-size: 20px;
  }

  .post-content {
    font-size: 13px;
  }

  .author-avatar {
    width: 32px;
    height: 32px;
    font-size: 14px;
  }

  .author-name {
    font-size: 13px;
  }

  .post-stats .stat-item {
    font-size: 12px;
  }

  .post-category {
    padding: 3px 8px;
    font-size: 11px;
  }
}
</style>
