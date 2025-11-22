import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../services/api';
import './TaskDetail.css';

export default function TaskDetail() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [task, setTask] = useState(null);
    const [progress, setProgress] = useState(0);
    const [notes, setNotes] = useState('');
    const [files, setFiles] = useState([]);
    const [uploading, setUploading] = useState(false);
    const [groupProgress, setGroupProgress] = useState([]);
    const [activeTab, setActiveTab] = useState('details');

    useEffect(() => {
        fetchTaskDetails();
        fetchProgress();
        fetchFiles();
        if (task?.is_group_task) {
            fetchGroupProgress();
        }
    }, [id]);

    const fetchTaskDetails = async () => {
        try {
            const response = await api.get(`/tasks/${id}`);
            setTask(response.data.data);
        } catch (error) {
            console.error('Error fetching task:', error);
            alert('خطا در دریافت اطلاعات تسک');
        }
    };

    const fetchProgress = async () => {
        try {
            const endpoint = task?.is_group_task ? `/tasks/${id}/my-progress` : `/tasks/${id}/progress`;
            const response = await api.get(endpoint);
            setProgress(response.data.data.progress || 0);
            setNotes(response.data.data.notes || '');
        } catch (error) {
            console.error('Error fetching progress:', error);
        }
    };

    const fetchFiles = async () => {
        try {
            const response = await api.get(`/tasks/${id}/files`);
            setFiles(response.data.data);
        } catch (error) {
            console.error('Error fetching files:', error);
        }
    };

    const fetchGroupProgress = async () => {
        try {
            const response = await api.get(`/tasks/${id}/progress`);
            setGroupProgress(response.data.data);
        } catch (error) {
            console.error('Error fetching group progress:', error);
        }
    };

    const handleProgressUpdate = async () => {
        try {
            const endpoint = task.is_group_task ? `/tasks/${id}/my-progress` : `/tasks/${id}/progress`;
            const data = task.is_group_task ?
                { progress: progress, notes: notes } :
                { progress: progress, notes: notes };

            await api.put(endpoint, data);
            alert('پیشرفت با موفقیت بروزرسانی شد');
            fetchProgress();
        } catch (error) {
            alert('خطا در بروزرسانی پیشرفت: ' + (error.response?.data?.error || 'Unknown error'));
        }
    };

    const handleFileUpload = async (event) => {
        const file = event.target.files[0];
        if (!file) return;

        const formData = new FormData();
        formData.append('file', file);
        formData.append('task_id', id);
        formData.append('notes', notes);

        setUploading(true);
        try {
            await api.post('/files/upload', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });
            alert('فایل با موفقیت آپلود شد');
            fetchFiles();
            event.target.value = ''; // Reset file input
        } catch (error) {
            alert('خطا در آپلود فایل: ' + (error.response?.data?.error || 'Unknown error'));
        } finally {
            setUploading(false);
        }
    };

    const handleDownloadFile = async (fileId, filename) => {
        try {
            const response = await api.get(`/files/${fileId}/download`, {
                responseType: 'blob'
            });

            const url = window.URL.createObjectURL(new Blob([response.data]));
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', filename);
            document.body.appendChild(link);
            link.click();
            link.remove();
            window.URL.revokeObjectURL(url);
        } catch (error) {
            alert('خطا در دانلود فایل');
        }
    };

    const handleDeleteFile = async (fileId) => {
        if (window.confirm('آیا از حذف این فایل مطمئن هستید؟')) {
            try {
                await api.delete(`/files/${fileId}`);
                alert('فایل با موفقیت حذف شد');
                fetchFiles();
            } catch (error) {
                alert('خطا در حذف فایل');
            }
        }
    };

    const updateMemberProgress = async (userId, newProgress) => {
        try {
            await api.put(`/tasks/${id}/group-progress`, {
                user_id: userId,
                progress: newProgress,
                approved: newProgress === 100
            });
            alert('پیشرفت عضو بروزرسانی شد');
            fetchGroupProgress();
        } catch (error) {
            alert('خطا در بروزرسانی پیشرفت عضو');
        }
    };

    if (!task) {
        return <div className="loading">در حال بارگذاری...</div>;
    }

    return (
        <div className="task-detail">
            <div className="container">
                <div className="task-header">
                    <button onClick={() => navigate(-1)} className="btn btn-secondary">
                        بازگشت
                    </button>
                    <h1>{task.title}</h1>
                    <div className={`task-status-badge status-${task.status}`}>
                        {task.status === 'pending' && 'در انتظار'}
                        {task.status === 'in_progress' && 'در حال انجام'}
                        {task.status === 'completed' && 'تکمیل شده'}
                        {task.status === 'expired' && 'منقضی شده'}
                    </div>
                </div>

                <div className="tabs">
                    <button
                        className={`tab-btn ${activeTab === 'details' ? 'active' : ''}`}
                        onClick={() => setActiveTab('details')}
                    >
                        جزئیات
                    </button>
                    <button
                        className={`tab-btn ${activeTab === 'progress' ? 'active' : ''}`}
                        onClick={() => setActiveTab('progress')}
                    >
                        پیشرفت
                    </button>
                    <button
                        className={`tab-btn ${activeTab === 'files' ? 'active' : ''}`}
                        onClick={() => setActiveTab('files')}
                    >
                        فایل‌ها ({files.length})
                    </button>
                    {task.is_group_task && (
                        <button
                            className={`tab-btn ${activeTab === 'members' ? 'active' : ''}`}
                            onClick={() => setActiveTab('members')}
                        >
                            وضعیت اعضا
                        </button>
                    )}
                </div>

                <div className="tab-content">
                    {/* تب جزئیات */}
                    {activeTab === 'details' && (
                        <div className="details-panel">
                            <div className="detail-item">
                                <label>توضیحات:</label>
                                <p>{task.description || 'بدون توضیحات'}</p>
                            </div>
                            <div className="detail-item">
                                <label>سازنده:</label>
                                <span>{task.creator?.username}</span>
                            </div>
                            <div className="detail-item">
                                <label>تاریخ ایجاد:</label>
                                <span>{new Date(task.created_at).toLocaleDateString('fa-IR')}</span>
                            </div>
                            {task.start_time && (
                                <div className="detail-item">
                                    <label>زمان شروع:</label>
                                    <span>{new Date(task.start_time).toLocaleDateString('fa-IR')}</span>
                                </div>
                            )}
                            {task.end_time && (
                                <div className="detail-item">
                                    <label>زمان پایان:</label>
                                    <span>{new Date(task.end_time).toLocaleDateString('fa-IR')}</span>
                                </div>
                            )}
                            {task.is_group_task && task.group && (
                                <div className="detail-item">
                                    <label>گروه:</label>
                                    <span>{task.group.name}</span>
                                </div>
                            )}
                        </div>
                    )}

                    {/* تب پیشرفت */}
                    {activeTab === 'progress' && (
                        <div className="progress-panel">
                            <div className="progress-section">
                                <h3>پیشرفت فعلی: {progress}%</h3>
                                <div className="progress-slider">
                                    <input
                                        type="range"
                                        min="0"
                                        max="100"
                                        value={progress}
                                        onChange={(e) => setProgress(parseInt(e.target.value))}
                                        className="slider"
                                    />
                                    <div className="slider-labels">
                                        <span>0%</span>
                                        <span>50%</span>
                                        <span>100%</span>
                                    </div>
                                </div>

                                <div className="progress-notes">
                                    <label>یادداشت‌ها:</label>
                                    <textarea
                                        value={notes}
                                        onChange={(e) => setNotes(e.target.value)}
                                        placeholder="یادداشت‌های خود را اینجا بنویسید..."
                                        rows="4"
                                    />
                                </div>

                                <button
                                    onClick={handleProgressUpdate}
                                    className="btn btn-primary"
                                    disabled={progress === (task.progress?.progress || 0)}
                                >
                                    ذخیره پیشرفت
                                </button>

                                {progress === 100 && (
                                    <div className="completion-message">
                                        ✅ تسک تکمیل شده است!
                                    </div>
                                )}
                            </div>
                        </div>
                    )}

                    {/* تب فایل‌ها */}
                    {activeTab === 'files' && (
                        <div className="files-panel">
                            <div className="upload-section">
                                <h3>آپلود فایل جدید</h3>
                                <div className="upload-area">
                                    <input
                                        type="file"
                                        onChange={handleFileUpload}
                                        disabled={uploading}
                                        className="file-input"
                                    />
                                    {uploading && <div className="uploading">در حال آپلود...</div>}
                                </div>
                            </div>

                            <div className="files-list">
                                <h3>فایل‌های آپلود شده ({files.length})</h3>
                                {files.length === 0 ? (
                                    <div className="empty-files">
                                        هنوز فایلی آپلود نشده است
                                    </div>
                                ) : (
                                    files.map((file) => (
                                        <div key={file.id} className="file-item">
                                            <div className="file-info">
                                                <span className="file-name">{file.filename}</span>
                                                <span className="file-size">{(file.file_size / 1024 / 1024).toFixed(2)} MB</span>
                                                <span className="file-uploader">آپلود شده توسط: {file.user?.username}</span>
                                            </div>
                                            <div className="file-actions">
                                                <button
                                                    onClick={() => handleDownloadFile(file.id, file.filename)}
                                                    className="btn btn-secondary btn-sm"
                                                >
                                                    دانلود
                                                </button>
                                                {(file.user_id === parseInt(localStorage.getItem('userId')) || task.creator_id === parseInt(localStorage.getItem('userId'))) && (
                                                    <button
                                                        onClick={() => handleDeleteFile(file.id)}
                                                        className="btn btn-danger btn-sm"
                                                    >
                                                        حذف
                                                    </button>
                                                )}
                                            </div>
                                        </div>
                                    ))
                                )}
                            </div>
                        </div>
                    )}

                    {/* تب وضعیت اعضا (فقط برای تسک‌های گروهی) */}
                    {activeTab === 'members' && task.is_group_task && (
                        <div className="members-progress-panel">
                            <h3>وضعیت پیشرفت اعضا</h3>
                            <div className="members-progress-list">
                                {groupProgress.map((memberProgress) => (
                                    <div key={memberProgress.id} className="member-progress-item">
                                        <div className="member-info">
                                            <span className="member-name">{memberProgress.user?.username}</span>
                                            <div className="progress-bar">
                                                <div
                                                    className="progress-fill"
                                                    style={{ width: `${memberProgress.progress}%` }}
                                                ></div>
                                                <span className="progress-text">{memberProgress.progress}%</span>
                                            </div>
                                        </div>
                                        <div className="member-actions">
                                            {task.creator_id === parseInt(localStorage.getItem('userId')) && (
                                                <input
                                                    type="range"
                                                    min="0"
                                                    max="100"
                                                    value={memberProgress.progress}
                                                    onChange={(e) => updateMemberProgress(memberProgress.user_id, parseInt(e.target.value))}
                                                    className="slider-small"
                                                />
                                            )}
                                            {memberProgress.approved && (
                                                <span className="approved-badge">تایید شده</span>
                                            )}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}