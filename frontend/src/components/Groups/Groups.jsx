import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import groupsService from '../../services/groups';
import './Groups.css';

function Groups() {
    const navigate = useNavigate();
    const [groups, setGroups] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [success, setSuccess] = useState(null);
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [formData, setFormData] = useState({
        name: '',
        description: '',
        members: []
    });
    const [creating, setCreating] = useState(false);
    const [searchQuery, setSearchQuery] = useState('');
    const [searchResults, setSearchResults] = useState([]);
    const [searching, setSearching] = useState(false);

    useEffect(() => {
        fetchGroups();
    }, []);

    const fetchGroups = async () => {
        try {
            setLoading(true);
            const response = await groupsService.getGroups();

            let groupsData = response.data?.data || response.data || [];

            if (groupsData && typeof groupsData === 'object' && !Array.isArray(groupsData)) {
                groupsData = Object.values(groupsData);
            }

            setGroups(groupsData);
        } catch (err) {
            setError('❌ خطا در دریافت گروه‌ها');
            console.error('❌ Fetch groups error:', err);
        } finally {
            setLoading(false);
        }
    };

    const handleSearchMembers = async (query) => {
        if (!query.trim()) {
            setSearchResults([]);
            return;
        }

        try {
            setSearching(true);
            const response = await groupsService.searchUsers(query);

            let results = response.data?.data || response.data || [];
            if (results && typeof results === 'object' && !Array.isArray(results)) {
                results = Object.values(results);
            }

            setSearchResults(results);
        } catch (err) {
            console.error('❌ Search error:', err);
            setSearchResults([]);
        } finally {
            setSearching(false);
        }
    };

    const handleAddMemberToForm = (user) => {
        const userId = user.id;

        if (!formData.members.find(m => m.id === userId)) {
            setFormData(prev => ({
                ...prev,
                members: [...prev.members, user]
            }));
            setSearchQuery('');
            setSearchResults([]);
        }
    };

    const handleRemoveMemberFromForm = (userId) => {
        setFormData(prev => ({
            ...prev,
            members: prev.members.filter(m => m.id !== userId)
        }));
    };

    const handleCreateGroup = async (e) => {
        e.preventDefault();
        if (!formData.name.trim()) {
            setError('❌ نام گروه ضروری است');
            return;
        }

        try {
            setCreating(true);
            setError(null);

            // ایجاد گروه
            const createResponse = await groupsService.createGroup({
                name: formData.name.trim(),
                description: formData.description.trim(),
                user_ids: formData.members.map(m => m.id)
            });

            setSuccess('✅ گروه با موفقیت ایجاد شد!');
            setFormData({ name: '', description: '', members: [] });
            setShowCreateForm(false);

            // منتظر و refresh
            setTimeout(() => {
                fetchGroups();
                setSuccess(null);
            }, 500);
        } catch (err) {
            setError(`❌ خطا در ایجاد گروه: ${err.response?.data?.error || err.message}`);
        } finally {
            setCreating(false);
        }
    };

    if (loading) {
        return <div className="loading">🔄 بارگذاری گروه‌ها...</div>;
    }

    return (
        <div className="page-container">
            <div className="page-header">
                <h1>👥 گروه‌های من</h1>
                <div className="header-actions">
                    <button
                        className="btn-primary"
                        onClick={() => setShowCreateForm(!showCreateForm)}
                    >
                        {showCreateForm ? '❌ انصراف' : '➕ گروه جدید'}
                    </button>
                    <button
                        className="btn-secondary"
                        onClick={() => navigate('/groups/invitations')}
                    >
                        📩 دعوت‌نامه‌ها
                    </button>
                </div>
            </div>

            {error && <div className="error-message">{error}</div>}
            {success && <div className="success-message">{success}</div>}

            {showCreateForm && (
                <div className="create-group-form">
                    <div className="form-header">
                        <h2>✨ ایجاد گروه جدید</h2>
                    </div>

                    <form onSubmit={handleCreateGroup} className="task-form">
                        <div className="form-group">
                            <label>📌 نام گروه *</label>
                            <input
                                type="text"
                                value={formData.name}
                                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                placeholder="نام گروه را وارد کنید"
                                required
                            />
                        </div>

                        <div className="form-group">
                            <label>📝 توضیحات</label>
                            <textarea
                                value={formData.description}
                                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                                placeholder="توضیحات گروه (اختیاری)"
                                rows="3"
                            />
                        </div>

                        <div className="form-group">
                            <label>👥 اضافه کردن اعضا (اختیاری)</label>
                            <div className="search-box">
                                <input
                                    type="text"
                                    placeholder="نام کاربری یا ایمیل را جستجو کنید..."
                                    value={searchQuery}
                                    onChange={(e) => {
                                        setSearchQuery(e.target.value);
                                        handleSearchMembers(e.target.value);
                                    }}
                                />
                                {searching && <div className="searching">🔍 در حال جستجو...</div>}
                            </div>

                            {searchResults.length > 0 && (
                                <div className="search-results">
                                    {searchResults.map(user => {
                                        const userId = user.id;
                                        const alreadyAdded = formData.members.find(m => m.id === userId);

                                        return (
                                            <div key={userId} className="search-result">
                                                <span>👤 {user.username} ({user.email})</span>
                                                <button
                                                    type="button"
                                                    className={alreadyAdded ? 'btn-remove' : 'btn-add'}
                                                    onClick={() => {
                                                        if (!alreadyAdded) {
                                                            handleAddMemberToForm(user);
                                                        }
                                                    }}
                                                    disabled={alreadyAdded}
                                                    style={{
                                                        opacity: alreadyAdded ? 0.6 : 1,
                                                        cursor: alreadyAdded ? 'not-allowed' : 'pointer'
                                                    }}
                                                >
                                                    {alreadyAdded ? '✅ اضافه شد' : '➕ اضافه'}
                                                </button>
                                            </div>
                                        );
                                    })}
                                </div>
                            )}
                        </div>

                        {formData.members.length > 0 && (
                            <div className="form-group">
                                <label>📋 اعضای منتخب ({formData.members.length})</label>
                                <div className="members-list">
                                    {formData.members.map(member => (
                                        <div key={member.id} className="member-checkbox">
                                            <label style={{ display: 'flex', alignItems: 'center' }}>
                                                ✅ {member.username} ({member.email})
                                            </label>
                                            <button
                                                type="button"
                                                className="btn-remove"
                                                onClick={() => handleRemoveMemberFromForm(member.id)}
                                                style={{ padding: '4px 8px', fontSize: '12px' }}
                                            >
                                                🗑️ حذف
                                            </button>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}

                        <div className="form-actions">
                            <button type="submit" className="btn-primary" disabled={creating}>
                                {creating ? '⏳ در حال ایجاد...' : '✅ ایجاد گروه'}
                            </button>
                            <button
                                type="button"
                                className="btn-secondary"
                                onClick={() => {
                                    setShowCreateForm(false);
                                    setFormData({ name: '', description: '', members: [] });
                                    setError(null);
                                }}
                            >
                                ❌ انصراف
                            </button>
                        </div>
                    </form>
                </div>
            )}

            {groups && groups.length > 0 ? (
                <div className="members-grid">
                    {groups.map((group) => {
                        const groupId = group.id;
                        const memberCount = group.members?.length || 0;

                        if (!groupId) {
                            return null;
                        }

                        return (
                            <div key={groupId} className="member-card">
                                <div className="member-header">
                                    <span className="member-name">📁 {group.name}</span>
                                    <span className="member-role admin">👑 مدیر</span>
                                </div>
                                <p className="member-email">{group.description || 'بدون توضیح'}</p>
                                <p style={{ fontSize: '12px', color: '#6b7280', marginBottom: '12px' }}>
                                    👥 اعضا: {memberCount}
                                </p>
                                <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap' }}>
                                    <button
                                        className="btn-primary"
                                        onClick={() => navigate(`/groups/${groupId}/add-task`)}
                                        style={{ flex: 1, minWidth: '100px' }}
                                    >
                                        ➕ تسک
                                    </button>
                                    <button
                                        className="btn-edit"
                                        onClick={() => navigate(`/groups/${groupId}/settings`)}
                                        style={{ flex: 1, minWidth: '100px' }}
                                    >
                                        ⚙️ تنظیمات
                                    </button>
                                </div>
                            </div>
                        );
                    })}
                </div>
            ) : (
                <div className="empty-state">
                    <p>📭 هیچ گروهی وجود ندارد</p>
                    <p className="empty-state-sub">
                        برای شروع، اولین گروه خود را ایجاد کنید
                    </p>
                    <button className="btn-primary" onClick={() => setShowCreateForm(true)}>
                        ➕ ایجاد گروه جدید
                    </button>
                </div>
            )}
        </div>
    );
}

export default Groups;