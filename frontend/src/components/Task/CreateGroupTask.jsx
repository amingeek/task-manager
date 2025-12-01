import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import groupsService from '../../services/groups';

function CreateGroupTask() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [groupLoading, setGroupLoading] = useState(true);
  const [error, setError] = useState(null);
  const [group, setGroup] = useState(null);
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    dueDate: '',
    priority: 'medium',
    assignedTo: [],
  });

  useEffect(() => {
    console.log('📍 Route params - id:', id);
    
    if (!id) {
      setError('❌ ID گروه نامعلوم است');
      setGroupLoading(false);
      return;
    }

    fetchGroup();
  }, [id]);

  const fetchGroup = async () => {
    try {
      setGroupLoading(true);
      console.log(`📥 Fetching group ${id}...`);
      const response = await groupsService.getGroup(id);
      console.log('✅ Group data:', response.data);
      
      const groupData = response.data?.data || response.data;
      setGroup(groupData);
      setError(null);
    } catch (err) {
      const msg = `❌ خطا در دریافت گروه: ${err.response?.status || err.message}`;
      setError(msg);
      console.error('❌ Fetch group error:', err);
    } finally {
      setGroupLoading(false);
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleMemberToggle = (userId) => {
    setFormData(prev => ({
      ...prev,
      assignedTo: prev.assignedTo.includes(userId)
        ? prev.assignedTo.filter(id => id !== userId)
        : [...prev.assignedTo, userId]
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      if (!formData.title.trim()) {
        setError('❌ عنوان تسک ضروری است');
        setLoading(false);
        return;
      }

      const taskData = {
        title: formData.title.trim(),
        description: formData.description.trim(),
        priority: formData.priority,
        status: 'pending',
      };

      if (formData.dueDate) {
        taskData.due_date = formData.dueDate;
      }

      console.log('📤 Sending task payload:', taskData);

      const response = await groupsService.createGroupTask(id, taskData);
      console.log('✅ Task created:', response.data);

      setError(null);
      navigate('/groups', { replace: true });
    } catch (err) {
      const msg = err.response?.data?.message || err.response?.data?.error || err.message || '❌ خطا در ایجاد تسک';
      setError(`❌ ${msg}`);
      console.error('❌ Error:', err);
    } finally {
      setLoading(false);
    }
  };

  if (groupLoading) {
    return <div className="loading">🔄 بارگذاری اطلاعات گروه...</div>;
  }

  if (!group) {
    return (
      <div className="task-detail-container">
        <div className="error-message">❌ {error || 'گروه پیدا نشد'}</div>
        <button className="btn-primary" onClick={() => navigate('/groups')}>
          ← برگشت به گروه‌ها
        </button>
      </div>
    );
  }

  return (
    <div className="task-detail-container">
      <div className="task-detail-header">
        <h1>✨ ایجاد تسک جدید</h1>
        <div style={{ display: 'flex', gap: '8px' }}>
          <button className="btn-back" onClick={() => navigate('/groups')}>← برگشت</button>
        </div>
      </div>

      <div className="group-info">
        <p><strong>📁 گروه:</strong> {group.name}</p>
        <p><strong>📝 توضیحات:</strong> {group.description || 'بدون توضیح'}</p>
      </div>

      {error && <div className="error-message">{error}</div>}

      <form onSubmit={handleSubmit} className="task-form">
        <div className="form-group">
          <label>📌 عنوان تسک *</label>
          <input
            type="text"
            name="title"
            value={formData.title}
            onChange={handleInputChange}
            placeholder="عنوان تسک را وارد کنید"
            required
          />
        </div>

        <div className="form-group">
          <label>📝 توضیحات</label>
          <textarea
            name="description"
            value={formData.description}
            onChange={handleInputChange}
            placeholder="توضیحات تسک (اختیاری)"
            rows="4"
          />
        </div>

        <div className="form-row">
          <div className="form-group">
            <label>📅 تاریخ سررسید</label>
            <input
              type="date"
              name="dueDate"
              value={formData.dueDate}
              onChange={handleInputChange}
            />
          </div>

          <div className="form-group">
            <label>⭐ اولویت</label>
            <select
              name="priority"
              value={formData.priority}
              onChange={handleInputChange}
            >
              <option value="low">🟢 کم</option>
              <option value="medium">🟡 متوسط</option>
              <option value="high">🔴 زیاد</option>
            </select>
          </div>
        </div>

        <div className="form-group">
          <label>👥 انتخاب اعضا (اختیاری)</label>
          <div className="members-list">
            {group.Members && group.Members.length > 0 ? (
              group.Members.map(member => (
                <div key={member.UserID || member.id} className="member-checkbox">
                  <input
                    type="checkbox"
                    id={`member-${member.UserID || member.id}`}
                    checked={formData.assignedTo.includes(member.UserID || member.id)}
                    onChange={() => handleMemberToggle(member.UserID || member.id)}
                  />
                  <label htmlFor={`member-${member.UserID || member.id}`}>
                    👤 {member.User?.username || member.username || 'کاربر'}
                    {member.Role === 'admin' && ' 👑'}
                  </label>
                </div>
              ))
            ) : (
              <p className="no-members">❌ هیچ عضوی در گروه نیست</p>
            )}
          </div>
        </div>

        <div className="form-actions">
          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? '⏳ در حال ایجاد...' : '✅ ایجاد تسک'}
          </button>
          <button type="button" className="btn-secondary" onClick={() => navigate('/groups')}>
            ❌ انصراف
          </button>
        </div>
      </form>
    </div>
  );
}

export default CreateGroupTask;
