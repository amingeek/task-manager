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
      
      if (response.data.data.token) {
        localStorage.setItem('token', response.data.data.token);
      }
      
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // ورود
  login: async (username, password) => {
    try {
      const response = await api.post('/login', {
        username,
        password,
      });
      
      if (response.data.data.token) {
        localStorage.setItem('token', response.data.data.token);
      }
      
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // خروج
  logout: () => {
    localStorage.removeItem('token');
  },

  // دریافت کاربر فعلی
  getCurrentUser: async () => {
    try {
      const response = await api.get('/me');
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // بروزرسانی پروفایل
  updateProfile: async (fullName, bio, avatarURL) => {
    try {
      const response = await api.put('/profile', {
        full_name: fullName,
        bio,
        avatar_url: avatarURL,
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // جستجو کاربران
  searchUsers: async (query) => {
    try {
      const response = await api.get(`/users/search?q=${query}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // دریافت streak
  getStreak: async () => {
    try {
      const response = await api.get('/analytics/streak');
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },
};

export default authService;
