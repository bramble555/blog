import { createRouter, createWebHistory } from 'vue-router'

// Layouts
import PortalLayout from '../layouts/PortalLayout.vue'
import AdminLayout from '../layouts/AdminLayout.vue'

// Portal Views
import HomeView from '../portal/HomeView.vue'
import ArticleDetail from '../portal/ArticleDetail.vue'
import AuthView from '../portal/AuthView.vue'

// Admin Views
import ArticleList from '../views/ArticleList.vue'
import ArticleEditor from '../views/ArticleEditor.vue'
import CollectionList from '../views/CollectionList.vue'
import UserList from '../views/UserList.vue'
import BannerList from '../views/BannerList.vue'
import AdvertList from '../views/AdvertList.vue'
import TagList from '../views/TagList.vue'
import CommentList from '../views/CommentList.vue'
import MessageList from '../views/MessageList.vue'
import ChatView from '../views/ChatView.vue'
import CalendarView from '../views/CalendarView.vue'
import Dashboard from '../views/Dashboard.vue'

const routes = [
    // Portal Routes
    {
        path: '/',
        component: PortalLayout,
        children: [
            {
                path: '',
                name: 'Home',
                component: HomeView,
                meta: { title: 'Home' }
            },
            {
                path: 'article/:sn',
                name: 'ArticleDetail',
                component: ArticleDetail,
                meta: { title: 'Article Detail' }
            }
        ]
    },
    // Auth Routes
    {
        path: '/login',
        name: 'Login',
        component: AuthView,
        meta: { title: 'Login' }
    },
    {
        path: '/register',
        name: 'Register',
        component: AuthView,
        meta: { title: 'Register' }
    },
    // Admin Routes
    {
        path: '/admin',
        component: AdminLayout,
        redirect: '/admin/dashboard',
        children: [
            {
                path: 'dashboard',
                name: 'AdminDashboard',
                component: Dashboard,
                meta: { title: 'Dashboard', requiresAuth: true, role: 2 }
            },
            {
                path: 'articles',
                name: 'AdminArticles',
                component: ArticleList,
                meta: { title: 'Article Management', requiresAuth: true, role: 2 }
            },
            {
                path: 'collections',
                name: 'AdminCollections',
                component: CollectionList,
                meta: { title: 'My Collections', requiresAuth: true, role: 2 }
            },
            {
                path: 'create',
                name: 'AdminArticleCreate',
                component: ArticleEditor,
                meta: { title: 'New Article', requiresAuth: true, role: 2 }
            },
            {
                path: 'edit/:sn',
                name: 'AdminArticleEdit',
                component: ArticleEditor,
                meta: { title: 'Edit Article', requiresAuth: true, role: 2 }
            },
            {
                path: 'users',
                name: 'AdminUsers',
                component: UserList,
                meta: { title: 'User Management', requiresAuth: true, role: 2 }
            },
            {
                path: 'banners',
                name: 'AdminBanners',
                component: BannerList,
                meta: { title: 'Banner Management', requiresAuth: true, role: 2 }
            },
            {
                path: 'adverts',
                name: 'AdminAdverts',
                component: AdvertList,
                meta: { title: 'Advert Management', requiresAuth: true, role: 2 }
            },
            {
                path: 'tags',
                name: 'AdminTags',
                component: TagList,
                meta: { title: 'Tag Management', requiresAuth: true, role: 2 }
            },
            {
                path: 'comments',
                name: 'AdminComments',
                component: CommentList,
                meta: { title: 'Comment Management', requiresAuth: true, role: 2 }
            },
            {
                path: 'calendar',
                name: 'AdminCalendar',
                component: CalendarView,
                meta: { title: 'Article Calendar', requiresAuth: true, role: 2 }
            },
            {
                path: 'messages',
                name: 'AdminMessages',
                component: MessageList,
                meta: { title: 'Messages', requiresAuth: true, role: 2 }
            },
            {
                path: 'chat',
                name: 'AdminChat',
                component: ChatView,
                meta: { title: 'System Chat', requiresAuth: true, role: 2 }
            }
        ]
    },
    // Fallback
    {
        path: '/:pathMatch(.*)*',
        redirect: '/'
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// Flag to ensure we only check for reload functionality once (on app init)
let hasCheckedReload = false

// Simple Route Guard
router.beforeEach((to, from, next) => {
    const title = to.meta?.title ? `${to.meta.title} | GVB Blog` : 'GVB Blog'
    document.title = title

    // Search Refresh Guard: 
    // Only run this check on the very first navigation (page load)
    if (!hasCheckedReload) {
        hasCheckedReload = true
        // If we are on the Home page with a search query, and the page was reloaded, redirect to '/'
        if (to.name === 'Home' && to.query.title) {
            const isReload = (
                (performance.getEntriesByType("navigation").length > 0 && performance.getEntriesByType("navigation")[0].type === 'reload') ||
                (window.performance && window.performance.navigation && window.performance.navigation.type === 1)
            );
            if (isReload) {
                next({ path: '/' });
                return;
            }
        }
    }

    // Auth Guard
    const token = localStorage.getItem('token')
    if (to.meta.requiresAuth && !token) {
        next('/login')
    } else {
        next()
    }
})

export default router
