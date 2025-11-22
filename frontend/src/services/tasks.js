// frontend/src/services/tasks.js

import api from './api';

export const tasksService = {
  // دریافت تمام تسکها
  getTasks: async () => {
    try {
      const response = await api.get('/tasks');
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // دریافت یک تسک
  getTask: async (id) => {
    try {
      const response = await api.get(`/tasks/${id}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // ایجاد تسک
  createTask: async (title, description, dueDate, startTime, endTime) => {
    try {
      const response = await api.post('/tasks', {
        title,
        description,
        due_date: dueDate,
        start_time: startTime,
        end_time: endTime,
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // بروزرسانی تسک
  updateTask: async (id, title, description, status) => {
    try {
      const response = await api.put(`/tasks/${id}`, {
        title,
        description,
        status,
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // حذف تسک
  deleteTask: async (id) => {
    try {
      const response = await api.delete(`/tasks/${id}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // بروزرسانی پیشرفت تسک شخصی
  updateProgress: async (id, progress, notes) => {
    try {
      const response = await api.put(`/tasks/${id}/progress`, {
        progress,
        notes,
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // دریافت پیشرفت تسک
  getProgress: async (id) => {
    try {
      const response = await api.get(`/tasks/${id}/progress`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // بروزرسانی پیشرفت تسک گروهی
  updateGroupProgress: async (id, userId, progress, approved, notes) => {
    try {
      const response = await api.put(`/tasks/${id}/group-progress`, {
        user_id: userId,
        progress,
        approved,
        notes,
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },
};

export default tasksService;
