import api from './api';

const groupsService = {
    // دریافت تمام گروه‌های کاربر
    getGroups: async () => {
        try {
            const response = await api.get('/groups');
            return response;
        } catch (error) {
            console.error('❌ getGroups error:', error);
            throw error;
        }
    },

    // دریافت یک گروه مشخص
    getGroup: async (id) => {
        try {
            const response = await api.get(`/groups/${id}`);
            return response;
        } catch (error) {
            console.error('❌ getGroup error:', error);
            throw error;
        }
    },

    // ایجاد گروه جدید
    createGroup: async (groupData) => {
        try {
            const response = await api.post('/groups', groupData);
            return response;
        } catch (error) {
            console.error('❌ createGroup error:', error);
            throw error;
        }
    },

    // به‌روزرسانی گروه
    updateGroup: async (id, groupData) => {
        try {
            const response = await api.put(`/groups/${id}`, groupData);
            return response;
        } catch (error) {
            console.error('❌ updateGroup error:', error);
            throw error;
        }
    },

    // جستجوی کاربران
    searchUsers: async (query) => {
        try {
            const response = await api.get(`/users/search?q=${encodeURIComponent(query)}`);
            return response;
        } catch (error) {
            console.error('❌ searchUsers error:', error);
            throw error;
        }
    },

    // اضافه کردن اعضا به گروه
    addMembers: async (groupId, memberData) => {
        try {
            const response = await api.post(`/groups/${groupId}/members`, memberData);
            return response;
        } catch (error) {
            console.error('❌ addMembers error:', error);
            throw error;
        }
    },

    // حذف یک عضو از گروه
    removeMember: async (groupId, userId) => {
        try {
            const response = await api.delete(`/groups/${groupId}/members/${userId}`);
            return response;
        } catch (error) {
            console.error('❌ removeMember error:', error);
            throw error;
        }
    },

    // حذف گروه
    deleteGroup: async (id) => {
        try {
            const response = await api.delete(`/groups/${id}`);
            return response;
        } catch (error) {
            console.error('❌ deleteGroup error:', error);
            throw error;
        }
    },

    // ایجاد تسک گروهی
    createGroupTask: async (groupId, taskData) => {
        try {
            const response = await api.post(`/groups/${groupId}/tasks`, taskData);
            return response;
        } catch (error) {
            console.error('❌ createGroupTask error:', error);
            throw error;
        }
    },

    // دریافت دعوت‌های در انتظار
    getPendingInvitations: async () => {
        try {
            const response = await api.get('/groups/invitations');
            return response;
        } catch (error) {
            console.error('❌ getPendingInvitations error:', error);
            throw error;
        }
    },

    // پذیرش دعوت گروه
    acceptInvitation: async (groupId) => {
        try {
            const response = await api.post(`/groups/${groupId}/accept`);
            return response;
        } catch (error) {
            console.error('❌ acceptInvitation error:', error);
            throw error;
        }
    }
};

export default groupsService;