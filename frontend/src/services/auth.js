// frontend/src/services/auth.js
import api from './api';

export const authService = {
    // ثبت نام
    register: async (username, email, password) => {
        try {
            const response = await api.post('/register', {
                username,
                email,
                password,
            });

            if (response.data.data && response.data.data.token) {
                localStorage.setItem('token', response.data.data.token);
            }

            return response;
        } catch (error) {
            throw error;
        }
    },

    // ورود
    login: async (username, password) => {
        try {
            const response = await api.post('/login', {
                username,
                password,
            });

            if (response.data.data && response.data.data.token) {
                localStorage.setItem('token', response.data.data.token);
            }

            return response;
        } catch (error) {
            throw error;
        }
    },

    // خروج
    logout: () => {
        localStorage.removeItem('token');
        localStorage.removeItem('userId');
    },

    // دریافت کاربر فعلی
    getCurrentUser: async () => {
        try {
            const response = await api.get('/me');
            return response;
        } catch (error) {
            throw error;
        }
    },

    // بروزرسانی پروفایل
    updateProfile: async (profileData) => {
        try {
            const response = await api.put('/profile', profileData);
            return response;
        } catch (error) {
            throw error;
        }
    },

    // جستجو کاربران
    searchUsers: async (query) => {
        try {
            const response = await api.get(`/users/search?q=${encodeURIComponent(query)}`);
            return response;
        } catch (error) {
            throw error;
        }
    },

    // دریافت streak
    getStreak: async () => {
        try {
            const response = await api.get('/analytics/streak');
            return response;
        } catch (error) {
            throw error;
        }
    },

    // دریافت آمار
    getAnalytics: async () => {
        try {
            const response = await api.get('/analytics/summary');
            return response;
        } catch (error) {
            throw error;
        }
    }
};

export default authService;