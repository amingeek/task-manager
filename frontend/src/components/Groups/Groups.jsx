import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../../services/api';
import './Groups.css';

export default function Groups() {
    const [groups, setGroups] = useState([]);
    const [searchQuery, setSearchQuery] = useState('');
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [newGroup, setNewGroup] = useState({ name: '', description: '', userIDs: '' });
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    useEffect(() => {
        fetchGroups();
    }, []);

    const fetchGroups = async () => {
        try {
            const response = await api.get('/groups');
            setGroups(response.data.data);
        } catch (error) {
            console.error('Error fetching groups:', error);
        }
    };

    const searchGroups = async () => {
        try {
            const response = await api.get(`/groups/search?q=${searchQuery}`);
            setGroups(response.data.data);
        } catch (error) {
            console.error('Error searching groups:', error);
        }
    };

    const handleCreateGroup = async (e) => {
        e.preventDefault();
        setLoading(true);
        try {
            const userIDs = newGroup.userIDs.split(',').map(id => parseInt(id.trim())).filter(id => id);
            await api.post('/groups', {
                name: newGroup.name,
                description: newGroup.description,
                userIDs: userIDs
            });
            setNewGroup({ name: '', description: '', userIDs: '' });
            setShowCreateForm(false);
            fetchGroups();
        } catch (error) {
            alert('خطا در ایجاد گروه: ' + (error.response?.data?.error || 'Unknown error'));
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="groups">
            <div className="container">
                <div className="groups-header">
                    <h1>مدیریت گروه‌ها</h1>
                    <div className="header-actions">
                        <button
                            onClick={() => setShowCreateForm(!showCreateForm)}
                            className="btn btn-primary"
                        >
                            ایجاد گروه جدید
                        </button>
                        <button onClick={() => navigate('/')} className="btn btn-secondary">
                            بازگشت به داشبورد
                        </button>
                    </div>
                </div>

                {/* فرم ایجاد گروه */}
                {showCreateForm && (
                    <div className="card create-group-form">
                        <h3>ایجاد گروه جدید</h3>
                        <form onSubmit={handleCreateGroup}>
                            <div className="form-group">
                                <label className="form-label">نام گروه</label>
                                <input
                                    type="text"
                                    value={newGroup.name}
                                    onChange={(e) => setNewGroup({...newGroup, name: e.target.value})}
                                    className="form-input"
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label className="form-label">توضیحات</label>
                                <textarea
                                    value={newGroup.description}
                                    onChange={(e) => setNewGroup({...newGroup, description: e.target.value})}
                                    className="form-input"
                                    rows="3"
                                />
                            </div>
                            <div className="form-group">
                                <label className="form-label">شناسه کاربران (با کاما جدا کنید)</label>
                                <input
                                    type="text"
                                    value={newGroup.userIDs}
                                    onChange={(e) => setNewGroup({...newGroup, userIDs: e.target.value})}
                                    className="form-input"
                                    placeholder="مثال: 1, 2, 3"
                                />
                            </div>
                            <div className="form-actions">
                                <button type="submit" className="btn btn-primary" disabled={loading}>
                                    {loading ? 'در حال ایجاد...' : 'ایجاد گروه'}
                                </button>
                                <button
                                    type="button"
                                    className="btn btn-secondary"
                                    onClick={() => setShowCreateForm(false)}
                                >
                                    انصراف
                                </button>
                            </div>
                        </form>
                    </div>
                )}

                {/* جستجو */}
                <div className="card search-section">
                    <div className="search-form">
                        <input
                            type="text"
                            placeholder="جستجو بین گروه‌ها..."
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            className="form-input"
                        />
                        <button onClick={searchGroups} className="btn btn-primary">
                            جستجو
                        </button>
                        <button onClick={() => { setSearchQuery(''); fetchGroups(); }} className="btn btn-secondary">
                            نمایش همه
                        </button>
                    </div>
                </div>

                {/* لیست گروه‌ها */}
                <div className="groups-list">
                    {groups.map((group) => (
                        <div key={group.id} className="card group-item">
                            <div className="group-header">
                                <div className="group-info">
                                    <h3>{group.name}</h3>
                                    <p className="group-description">{group.description}</p>
                                    <div className="group-meta">
                                        <span>سازنده: {group.creator?.username}</span>
                                        <span>اعضا: {group.members?.length} نفر</span>
                                        <span>تسک‌ها: {group.tasks?.length} مورد</span>
                                    </div>
                                </div>
                                <div className="group-actions">
                                    <button
                                        onClick={() => navigate(`/groups/${group.id}/tasks`)}
                                        className="btn btn-primary"
                                    >
                                        مشاهده تسک‌ها
                                    </button>
                                    <button
                                        onClick={() => navigate(`/groups/${group.id}/add-task`)}
                                        className="btn btn-secondary"
                                    >
                                        افزودن تسک
                                    </button>
                                </div>
                            </div>

                            {/* اعضای گروه */}
                            <div className="group-members">
                                <h4>اعضای گروه:</h4>
                                <div className="members-list">
                                    {group.members?.map((member) => (
                                        <div key={member.id} className="member-item">
                                            <span className="member-name">{member.user?.username}</span>
                                            <span className={`member-role ${member.role}`}>
                        {member.role === 'admin' ? 'مدیر' : 'عضو'}
                      </span>
                                            <span className={`member-status ${member.accepted ? 'accepted' : 'pending'}`}>
                        {member.accepted ? 'تایید شده' : 'در انتظار'}
                      </span>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        </div>
                    ))}
                </div>

                {groups.length === 0 && (
                    <div className="empty-state">
                        <p>هنوز گروهی وجود ندارد!</p>
                        <p className="empty-state-sub">اولین گروه خود را ایجاد کنید.</p>
                    </div>
                )}
            </div>
        </div>
    );
}