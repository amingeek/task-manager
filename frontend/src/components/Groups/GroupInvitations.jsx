import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import groupsService from '../../services/groups';
import './Groups.css';

function GroupInvitations() {
    const [invitations, setInvitations] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [success, setSuccess] = useState(null);
    const navigate = useNavigate();

    useEffect(() => {
        fetchPendingInvitations();
    }, []);

    const fetchPendingInvitations = async () => {
        try {
            setLoading(true);
            const response = await groupsService.getPendingInvitations();
            setInvitations(response.data.data || []);
            setError(null);
        } catch (err) {
            setError('❌ خطا در دریافت دعوت‌ها');
            console.error('❌ Fetch invitations error:', err);
        } finally {
            setLoading(false);
        }
    };

    const handleAcceptInvitation = async (groupId) => {
        try {
            setError(null);
            await groupsService.acceptInvitation(groupId);
            setSuccess('✅ دعوت با موفقیت پذیرفته شد!');

            // به روزرسانی لیست دعوت‌ها
            setTimeout(() => {
                fetchPendingInvitations();
                setSuccess(null);
            }, 2000);
        } catch (err) {
            setError('❌ خطا در پذیرش دعوت');
            console.error('❌ Accept invitation error:', err);
        }
    };

    const handleRejectInvitation = async (groupId) => {
        if (!window.confirm('آیا از رد این دعوت مطمئن هستید؟')) return;

        try {
            // اگر endpoint برای رد دعوت دارید، اینجا اضافه کنید
            // در حال حاضر فقط از لیست حذف می‌کنیم
            setInvitations(prev => prev.filter(inv => inv.group_id !== groupId));
            setSuccess('✅ دعوت رد شد');
            setTimeout(() => setSuccess(null), 3000);
        } catch (err) {
            setError('❌ خطا در رد دعوت');
        }
    };

    if (loading) {
        return <div className="loading">🔄 بارگذاری دعوت‌ها...</div>;
    }

    return (
        <div className="page-container">
            <div className="page-header">
                <h1>📩 دعوت‌نامه‌های گروه</h1>
                <button className="btn-back" onClick={() => navigate('/groups')}>
                    ← بازگشت به گروه‌ها
                </button>
            </div>

            {error && <div className="error-message">{error}</div>}
            {success && <div className="success-message">{success}</div>}

            {invitations.length === 0 ? (
                <div className="empty-state">
                    <div className="empty-icon">📭</div>
                    <p>هیچ دعوت‌نامه‌ای در انتظار ندارید</p>
                    <p className="empty-state-sub">
                        وقتی مدیران گروه شما را به گروهی دعوت کنند، اینجا نمایش داده می‌شود.
                    </p>
                </div>
            ) : (
                <div className="members-grid">
                    {invitations.map((invitation) => {
                        const group = invitation.Group;

                        return (
                            <div key={invitation.id} className="member-card">
                                <div className="member-header">
                                    <span className="member-name">📁 {group.name}</span>
                                    <span className="member-role pending">⏳ در انتظار</span>
                                </div>

                                <p className="member-email">
                                    {group.description || 'بدون توضیح'}
                                </p>

                                <div className="invitation-details">
                                    <p style={{ fontSize: '12px', color: '#6b7280', marginBottom: '12px' }}>
                                        👤 دعوت شده توسط: {group.Creator?.username || 'نامعلوم'}
                                    </p>
                                    <p style={{ fontSize: '12px', color: '#6b7280' }}>
                                        📅 تاریخ دعوت: {new Date(invitation.created_at).toLocaleDateString('fa-IR')}
                                    </p>
                                </div>

                                <div style={{ display: 'flex', gap: '8px', marginTop: '16px' }}>
                                    <button
                                        className="btn-primary"
                                        onClick={() => handleAcceptInvitation(group.id)}
                                        style={{ flex: 1 }}
                                    >
                                        ✅ پذیرش دعوت
                                    </button>
                                    <button
                                        className="btn-remove"
                                        onClick={() => handleRejectInvitation(group.id)}
                                        style={{ flex: 1 }}
                                    >
                                        ❌ رد دعوت
                                    </button>
                                </div>
                            </div>
                        );
                    })}
                </div>
            )}
        </div>
    );
}

export default GroupInvitations;