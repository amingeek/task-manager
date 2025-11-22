// frontend/src/services/groups.js

import api from './api';

export const groupsService = {
  // دریافت تمام گروهها
  getGroups: async () => {
    try {
      const response = await api.get('/groups');
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // دریافت یک گروه
  getGroup: async (id) => {
    try {
      const response = await api.get(`/groups/${id}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // ایجاد گروه
  createGroup: async (name, description, userIds = []) => {
    try {
      const response = await api.post('/groups', {
        name,
        description,
        user_ids: userIds,
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // بروزرسانی گروه
  updateGroup: async (id, name, description) => {
    try {
      const response = await api.put(`/groups/${id}`, {
        name,
        description,
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // حذف گروه
  deleteGroup: async (id) => {
    try {
      const response = await api.delete(`/groups/${id}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // جستجو گروهها
  searchGroups: async (query) => {
    try {
      const response = await api.get(`/groups/search?q=${query}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // اضافه کردن اعضا به گروه
  addMembers: async (groupId, userIds) => {
    try {
      const response = await api.post(`/groups/${groupId}/members`, {
        user_ids: userIds,
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // حذف عضو از گروه
  removeMember: async (groupId, userId) => {
    try {
      const response = await api.delete(`/groups/${groupId}/members/${userId}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // پذیرفتن دعوت گروه
  acceptInvitation: async (groupId, userId) => {
    try {
      const response = await api.post(`/groups/${groupId}/members/${userId}/accept`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // ایجاد تسک گروهی
  createGroupTask: async (groupId, title, description, dueDate, userIds = []) => {
    try {
      const response = await api.post(`/groups/${groupId}/tasks`, {
        title,
        description,
        due_date: dueDate,
        user_ids: userIds,
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // دریافت تسکهای گروه
  getGroupTasks: async (groupId) => {
    try {
      const response = await api.get(`/groups/${groupId}/tasks`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // بروزرسانی تسک گروهی
  updateGroupTask: async (groupId, taskId, title, description, status) => {
    try {
      const response = await api.put(`/groups/${groupId}/tasks/${taskId}`, {
        title,
        description,
        status,
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },

  // حذف تسک گروهی
  deleteGroupTask: async (groupId, taskId) => {
    try {
      const response = await api.delete(`/groups/${groupId}/tasks/${taskId}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error;
    }
  },
};

export default groupsService;
