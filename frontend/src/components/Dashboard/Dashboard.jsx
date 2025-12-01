// frontend/src/components/Dashboard/Dashboard.jsx
import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../../services/api';
import './Dashboard.css';

export default function Dashboard() {
    const [user, setUser] = useState(null);
    const [tasks, setTasks] = useState([]);
    const [groups, setGroups] = useState([]);
    const [newTaskTitle, setNewTaskTitle] = useState('');
    const [newTaskDescription, setNewTaskDescription] = useState('');
    const [newTaskStartTime, setNewTaskStartTime] = useState('');
    const [newTaskEndTime, setNewTaskEndTime] = useState('');
    const [activeTab, setActiveTab] = useState('personal');
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            navigate('/login');
            return;
        }

        fetchUserData();
        fetchPersonalTasks();
        fetchUserGroups();
    }, [navigate]);

    const fetchUserData = async () => {
        try {
            const response = await api.get('/me');
            setUser(response.data.data);
        } catch (error) {
            console.error('Error fetching user:', error);
            localStorage.removeItem('token');
            navigate('/login');
        }
    };

    const fetchPersonalTasks = async () => {
        try {
            const response = await api.get('/tasks');
            setTasks(response.data.data || []);
        } catch (error) {
            console.error('Error fetching tasks:', error);
            setTasks([]);
        }
    };

    const fetchUserGroups = async () => {
        try {
            const response = await api.get('/groups');
            setGroups(response.data.data || []);
        } catch (error) {
            console.error('Error fetching groups:', error);
            setGroups([]);
        }
    };

    const handleCreateTask = async (e) => {
        e.preventDefault();
        if (!newTaskTitle.trim()) return;

        setLoading(true);
        try {
            const taskData = {
                title: newTaskTitle,
                description: newTaskDescription
            };

            // Add start time if provided
            if (newTaskStartTime) {
                taskData.start_time = new Date(newTaskStartTime).toISOString();
            }

            // Add end time if provided
            if (newTaskEndTime) {
                taskData.end_time = new Date(newTaskEndTime).toISOString();
            }

            await api.post('/tasks', taskData);
            setNewTaskTitle('');
            setNewTaskDescription('');
            setNewTaskStartTime('');
            setNewTaskEndTime('');
            fetchPersonalTasks();
        } catch (error) {
            alert('خطا در ایجاد تسک: ' + (error.response?.data?.error || 'Unknown error'));
        } finally {
            setLoading(false);
        }
    };

    const handleUpdateTaskStatus = async (taskId, newStatus) => {
        try {
            await api.put(`/tasks/${taskId}`, { status: newStatus });
            fetchPersonalTasks();
        } catch (error) {
            alert('خطا در بروزرسانی تسک: ' + (error.response?.data?.error || 'Unknown error'));
        }
    };

    const handleDeleteTask = async (taskId) => {
        if (window.confirm('آیا از حذف این تسک مطمئن هستید؟')) {
            try {
                await api.delete(`/tasks/${taskId}`);
                fetchPersonalTasks();
            } catch (error) {
                alert('خطا در حذف تسک: ' + (error.response?.data?.error || 'Unknown error'));
            }
        }
    };

    const handleLogout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('userId');
        navigate('/login');
    };

    const getTaskStatusColor = (status) => {
        switch (status) {
            case 'completed': return '#28a745';
            case 'expired': return '#dc3545';
            default: return '#ffc107';
        }
    };

    const getTaskStatusText = (status) => {
        switch (status) {
            case 'completed': return 'انجام شده';
            case 'expired': return 'منقضی شده';
            default: return 'در انتظار';
        }
    };

    const formatDate = (dateString) => {
        if (!dateString) return '-';
        return new Date(dateString).toLocaleDateString('fa-IR', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    // آمارها
    const personalTasks = tasks.filter(task => !task.is_group_task);
    const completedTasks = personalTasks.filter(task => task.status === 'completed').length;
    const pendingTasks = personalTasks.filter(task => task.status === 'pending').length;
    const expiredTasks = personalTasks.filter(task => task.status === 'expired').length;

    return (
        <div className="dashboard">
            {/* هدر */}
            <header className="dashboard-header">
                <div className="container">
                    <div className="header-content">
                        <div className="header-info">
                            <h1>تسک منیجر</h1>
                            {user && (
                                <div className="user-info">
                                    <span>خوش آمدید، <strong>{user.username}</strong></span>
                                    <span className="user-email">{user.email}</span>
                                    <span className="user-id">ID: {user.id}</span>
                                </div>
                            )}
                        </div>
                        <div className="header-actions">
                            <button onClick={() => navigate('/profile')} className="btn btn-secondary">
                                پروفایل
                            </button>
                            <button onClick={handleLogout} className="btn btn-secondary logout-btn">
                                خروج
                            </button>
                        </div>
                    </div>
                </div>
            </header>

            {/* محتوای اصلی */}
            <main className="dashboard-main">
                <div className="container">
                    {/* کارت‌های آمار */}
                    <div className="stats-section">
                        <div className="stats-grid">
                            <div className="stat-card total-tasks">
                                <div className="stat-icon">📝</div>
                                <div className="stat-content">
                                    <div className="stat-number">{personalTasks.length}</div>
                                    <div className="stat-label">کل تسک‌ها</div>
                                </div>
                            </div>

                            <div className="stat-card completed-tasks">
                                <div className="stat-icon">✅</div>
                                <div className="stat-content">
                                    <div className="stat-number">{completedTasks}</div>
                                    <div className="stat-label">انجام شده</div>
                                </div>
                            </div>

                            <div className="stat-card pending-tasks">
                                <div className="stat-icon">⏳</div>
                                <div className="stat-content">
                                    <div className="stat-number">{pendingTasks}</div>
                                    <div className="stat-label">در انتظار</div>
                                </div>
                            </div>

                            <div className="stat-card groups-count">
                                <div className="stat-icon">👥</div>
                                <div className="stat-content">
                                    <div className="stat-number">{groups.length}</div>
                                    <div className="stat-label">گروه‌ها</div>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* تب‌های اصلی */}
                    <div className="tabs-section">
                        <div className="tabs-header">
                            <button
                                className={`tab-btn ${activeTab === 'personal' ? 'active' : ''}`}
                                onClick={() => setActiveTab('personal')}
                            >
                                تسک‌های شخصی
                            </button>
                            <button
                                className={`tab-btn ${activeTab === 'groups' ? 'active' : ''}`}
                                onClick={() => setActiveTab('groups')}
                            >
                                گروه‌های من
                            </button>
                            <button
                                className={`tab-btn ${activeTab === 'create' ? 'active' : ''}`}
                                onClick={() => setActiveTab('create')}
                            >
                                ایجاد تسک جدید
                            </button>
                        </div>

                        <div className="tab-content">
                            {/* تب تسک‌های شخصی */}
                            {activeTab === 'personal' && (
                                <div className="tab-panel">
                                    <div className="section-header">
                                        <h3>تسک‌های شخصی شما</h3>
                                        <span className="tasks-count">({personalTasks.length} تسک)</span>
                                    </div>

                                    {personalTasks.length === 0 ? (
                                        <div className="empty-state">
                                            <div className="empty-icon">📝</div>
                                            <p>هنوز تسکی ایجاد نکرده‌اید!</p>
                                            <p className="empty-state-sub">برای ایجاد اولین تسک، به تب "ایجاد تسک جدید" بروید.</p>
                                        </div>
                                    ) : (
                                        <div className="tasks-list">
                                            {personalTasks.map((task) => (
                                                <div key={task.id} className="task-item">
                                                    <div className="task-main">
                                                        <div className="task-header">
                                                            <h4 className="task-title">{task.title}</h4>
                                                            <div className="task-actions">
                                                                {task.status === 'pending' && (
                                                                    <button
                                                                        onClick={() => handleUpdateTaskStatus(task.id, 'completed')}
                                                                        className="btn btn-success btn-sm"
                                                                    >
                                                                        انجام شد
                                                                    </button>
                                                                )}
                                                                {task.status === 'completed' && (
                                                                    <button
                                                                        onClick={() => handleUpdateTaskStatus(task.id, 'pending')}
                                                                        className="btn btn-warning btn-sm"
                                                                    >
                                                                        بازگشت به انتظار
                                                                    </button>
                                                                )}
                                                                <button
                                                                    onClick={() => handleDeleteTask(task.id)}
                                                                    className="btn btn-danger btn-sm"
                                                                >
                                                                    حذف
                                                                </button>
                                                            </div>
                                                        </div>

                                                        {task.description && (
                                                            <p className="task-description">{task.description}</p>
                                                        )}

                                                        <div className="task-meta">
                                                            <div className="meta-item">
                                                                <span className="meta-label">وضعیت:</span>
                                                                <span
                                                                    className="task-status"
                                                                    style={{ color: getTaskStatusColor(task.status) }}
                                                                >
                                  {getTaskStatusText(task.status)}
                                </span>
                                                            </div>

                                                            <div className="meta-item">
                                                                <span className="meta-label">تاریخ ایجاد:</span>
                                                                <span className="meta-value">{formatDate(task.created_at)}</span>
                                                            </div>

                                                            {task.start_time && (
                                                                <div className="meta-item">
                                                                    <span className="meta-label">زمان شروع:</span>
                                                                    <span className="meta-value">{formatDate(task.start_time)}</span>
                                                                </div>
                                                            )}

                                                            {task.end_time && (
                                                                <div className="meta-item">
                                                                    <span className="meta-label">زمان پایان:</span>
                                                                    <span className="meta-value">{formatDate(task.end_time)}</span>
                                                                </div>
                                                            )}
                                                        </div>
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    )}
                                </div>
                            )}

                            {/* تب گروه‌ها */}
                            {activeTab === 'groups' && (
                                <div className="tab-panel">
                                    <div className="section-header">
                                        <h3>گروه‌های شما</h3>
                                        <button
                                            onClick={() => navigate('/groups')}
                                            className="btn btn-primary"
                                        >
                                            مدیریت گروه‌ها
                                        </button>
                                    </div>

                                    {groups.length === 0 ? (
                                        <div className="empty-state">
                                            <div className="empty-icon">👥</div>
                                            <p>هنوز در گروهی عضو نیستید!</p>
                                            <p className="empty-state-sub">
                                                <button
                                                    onClick={() => navigate('/groups')}
                                                    className="btn btn-primary"
                                                >
                                                    اولین گروه خود را ایجاد کنید
                                                </button>
                                            </p>
                                        </div>
                                    ) : (
                                        <div className="groups-grid">
                                            {groups.map((group) => (
                                                <div key={group.id} className="group-card">
                                                    <div className="group-header">
                                                        <h4 className="group-name">{group.name}</h4>
                                                        <span className="group-task-count">
                              {group.tasks?.length || 0} تسک
                            </span>
                                                    </div>

                                                    {group.description && (
                                                        <p className="group-description">{group.description}</p>
                                                    )}

                                                    <div className="group-meta">
                                                        <div className="meta-item">
                                                            <span className="meta-label">سازنده:</span>
                                                            <span className="meta-value">{group.creator?.username}</span>
                                                        </div>
                                                        <div className="meta-item">
                                                            <span className="meta-label">اعضا:</span>
                                                            <span className="meta-value">{group.members?.length} نفر</span>
                                                        </div>
                                                    </div>

                                                    <div className="group-actions">
                                                        <button
                                                            onClick={() => navigate(`/groups/${group.id}/tasks`)}
                                                            className="btn btn-primary btn-sm"
                                                        >
                                                            مشاهده تسک‌ها
                                                        </button>
                                                        <button
                                                            onClick={() => navigate(`/groups/${group.id}/add-task`)}
                                                            className="btn btn-secondary btn-sm"
                                                        >
                                                            افزودن تسک
                                                        </button>
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    )}
                                </div>
                            )}

                            {/* تب ایجاد تسک جدید */}
                            {activeTab === 'create' && (
                                <div className="tab-panel">
                                    <div className="section-header">
                                        <h3>ایجاد تسک جدید</h3>
                                    </div>

                                    <div className="create-task-form">
                                        <form onSubmit={handleCreateTask}>
                                            <div className="form-group">
                                                <label className="form-label">عنوان تسک *</label>
                                                <input
                                                    type="text"
                                                    placeholder="عنوان تسک را وارد کنید..."
                                                    value={newTaskTitle}
                                                    onChange={(e) => setNewTaskTitle(e.target.value)}
                                                    className="form-input"
                                                    required
                                                />
                                            </div>

                                            <div className="form-group">
                                                <label className="form-label">توضیحات</label>
                                                <textarea
                                                    placeholder="توضیحات تسک (اختیاری)..."
                                                    value={newTaskDescription}
                                                    onChange={(e) => setNewTaskDescription(e.target.value)}
                                                    className="form-input"
                                                    rows="4"
                                                />
                                            </div>

                                            <div className="form-row">
                                                <div className="form-group">
                                                    <label className="form-label">زمان شروع</label>
                                                    <input
                                                        type="datetime-local"
                                                        value={newTaskStartTime}
                                                        onChange={(e) => setNewTaskStartTime(e.target.value)}
                                                        className="form-input"
                                                    />
                                                </div>

                                                <div className="form-group">
                                                    <label className="form-label">زمان پایان</label>
                                                    <input
                                                        type="datetime-local"
                                                        value={newTaskEndTime}
                                                        onChange={(e) => setNewTaskEndTime(e.target.value)}
                                                        className="form-input"
                                                    />
                                                </div>
                                            </div>

                                            <div className="form-actions">
                                                <button
                                                    type="submit"
                                                    className="btn btn-primary"
                                                    disabled={loading || !newTaskTitle.trim()}
                                                >
                                                    {loading ? 'در حال ایجاد...' : 'ایجاد تسک'}
                                                </button>
                                                <button
                                                    type="button"
                                                    className="btn btn-secondary"
                                                    onClick={() => {
                                                        setNewTaskTitle('');
                                                        setNewTaskDescription('');
                                                        setNewTaskStartTime('');
                                                        setNewTaskEndTime('');
                                                    }}
                                                >
                                                    پاک کردن فرم
                                                </button>
                                            </div>
                                        </form>
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>

                    {/* منوی سریع */}
                    <div className="quick-actions">
                        <h3>دسترسی سریع</h3>
                        <div className="actions-grid">
                            <button onClick={() => navigate('/profile')} className="action-btn profile-btn">
                                <span className="action-icon">👤</span>
                                <span className="action-text">پروفایل کاربری</span>
                            </button>

                            <button onClick={() => navigate('/groups')} className="action-btn groups-btn">
                                <span className="action-icon">👥</span>
                                <span className="action-text">مدیریت گروه‌ها</span>
                            </button>

                            <button onClick={() => setActiveTab('create')} className="action-btn create-task-btn">
                                <span className="action-icon">➕</span>
                                <span className="action-text">تسک جدید</span>
                            </button>

                            <button onClick={() => setActiveTab('personal')} className="action-btn tasks-btn">
                                <span className="action-icon">📋</span>
                                <span className="action-text">مشاهده تسک‌ها</span>
                            </button>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    );
}