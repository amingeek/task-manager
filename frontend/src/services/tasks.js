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

  // ایجاد تسک شخصی
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
  updateProgress: async (id, progress, notes = '') => {
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

  // دریافت پیشرفت تسک شخصی
  getProgress: async (id) => {
    try {
      const response = await api.get(`/tasks/${id}/progress`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // دریافت پیشرفت شخصی برای تسک گروهی
  getMyGroupProgress: async (id) => {
    try {
      const response = await api.get(`/tasks/${id}/progress`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // آپلود فایل برای تسک
  uploadFile: async (taskId, file, notes = '') => {
    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('task_id', taskId);
      formData.append('notes', notes);

      const response = await api.post(`/tasks/${taskId}/files`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // دریافت فایل‌های تسک
  getTaskFiles: async (taskId) => {
    try {
      const response = await api.get(`/tasks/${taskId}/files`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // دانلود فایل
  downloadFile: async (fileId) => {
    try {
      const response = await api.get(`/files/${fileId}`, {
        responseType: 'blob',
      });
      return response;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // حذف فایل
  deleteFile: async (fileId) => {
    try {
      const response = await api.delete(`/files/${fileId}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },
};

export default tasksService;
