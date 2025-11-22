import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../../services/api';
import './Profile.css';

export default function Profile() {
    const [user, setUser] = useState(null);
    const [stats, setStats] = useState({});
    const navigate = useNavigate();

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            navigate('/login');
            return;
        }
        fetchUserData();
        fetchUserStats();
    }, [navigate]);

    const fetchUserData = async () => {
        try {
            const response = await api.get('/me');
            setUser(response.data.data);
        } catch (error) {
            console.error('Error fetching user:', error);
        }
    };

    const fetchUserStats = async () => {
        try {
            const tasksResponse = await api.get('/tasks');
            const tasks = tasksResponse.data.data;

            const personalTasks = tasks.filter(task => !task.is_group_task);
            const completedTasks = personalTasks.filter(task => task.status === 'completed').length;

            setStats({
                totalTasks: personalTasks.length,
                completedTasks: completedTasks,
                completionRate: personalTasks.length > 0 ? Math.round((completedTasks / personalTasks.length) * 100) : 0
            });
        } catch (error) {
            console.error('Error fetching stats:', error);
        }
    };

    if (!user) {
        return <div>Loading...</div>;
    }

    return (
        <div className="profile">
            <div className="container">
                <div className="profile-header">
                    <h1>پروفایل کاربری</h1>
                    <button onClick={() => navigate('/')} className="btn btn-secondary">
                        بازگشت به داشبورد
                    </button>
                </div>

                <div className="profile-grid">
                    {/* اطلاعات کاربر */}
                    <div className="card profile-card">
                        <h3>اطلاعات شخصی</h3>
                        <div className="profile-info">
                            <div className="info-item">
                                <label>شناسه کاربری:</label>
                                <span className="user-id">{user.id}</span>
                            </div>
                            <div className="info-item">
                                <label>نام کاربری:</label>
                                <span>{user.username}</span>
                            </div>
                            <div className="info-item">
                                <label>ایمیل:</label>
                                <span>{user.email}</span>
                            </div>
                            <div className="info-item">
                                <label>تاریخ عضویت:</label>
                                <span>{new Date(user.created_at).toLocaleDateString('fa-IR')}</span>
                            </div>
                        </div>
                    </div>

                    {/* آمار کاربر */}
                    <div className="card stats-card">
                        <h3>آمار عملکرد</h3>
                        <div className="stats-grid">
                            <div className="stat-item">
                                <div className="stat-number">{stats.totalTasks}</div>
                                <div className="stat-label">تسک‌های شخصی</div>
                            </div>
                            <div className="stat-item">
                                <div className="stat-number">{stats.completedTasks}</div>
                                <div className="stat-label">تسک‌های انجام شده</div>
                            </div>
                            <div className="stat-item">
                                <div className="stat-number">{stats.completionRate}%</div>
                                <div className="stat-label">نرخ تکمیل</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}