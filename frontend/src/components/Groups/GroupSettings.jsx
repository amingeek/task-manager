import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import groupsService from '../../services/groups';

function GroupSettings() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const [group, setGroup] = useState(null);
  const [editMode, setEditMode] = useState(false);
  const [editName, setEditName] = useState('');
  const [editDescription, setEditDescription] = useState('');
  const [newMemberUsername, setNewMemberUsername] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [searching, setSearching] = useState(false);

  useEffect(() => {
    if (!id) {
      setError('❌ ID گروه نامعلوم است');
      setLoading(false);
      return;
    }
    fetchGroup();
  }, [id]);

  const fetchGroup = async () => {
    try {
      setLoading(true);
      const response = await groupsService.getGroup(id);
      
      const groupData = response.data?.data || response.data;
      setGroup(groupData);
      setEditName(groupData.name);
      setEditDescription(groupData.description || '');
      setError(null);
    } catch (err) {
      setError(`❌ خطا در دریافت گروه: ${err.response?.status || err.message}`);
      console.error('❌ Fetch group error:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateGroup = async (e) => {
    e.preventDefault();
    try {
      setLoading(true);
      await groupsService.updateGroup(id, {
        name: editName.trim(),
        description: editDescription.trim(),
      });
      setSuccess('✅ گروه با موفقیت به‌روزرسانی شد');
      setEditMode(false);
      fetchGroup();
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('❌ خطا در به‌روزرسانی گروه');
      console.error(err);
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
      const results = response.data?.data || response.data || [];
      
      // فیلتر کریں - تنہا افراد جو موجود نہیں
      const existingIds = (group?.members || []).map(m => m.user_id || m.id);
      const filtered = results.filter(user => !existingIds.includes(user.id));
      
      setSearchResults(filtered);
    } catch (err) {
      console.error('❌ Search error:', err);
      setSearchResults([]);
    } finally {
      setSearching(false);
    }
  };

  const handleAddMember = async (userId) => {
    try {
      await groupsService.addMembers(id, { user_ids: [userId] });
      setSuccess('✅ عضو با موفقیت اضافه شد!');
      setNewMemberUsername('');
      setSearchResults([]);
      fetchGroup();
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('❌ خطا در اضافه کردن عضو');
      console.error(err);
    }
  };

  const handleRemoveMember = async (userId) => {
    if (!window.confirm('آیا مطمئن هستید؟')) return;

    try {
      await groupsService.removeMember(id, userId);
      setSuccess('✅ عضو با موفقیت حذف شد');
      fetchGroup();
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError('❌ خطا در حذف عضو');
      console.error(err);
    }
  };

  const handleDeleteGroup = async () => {
    if (!window.confirm('آیا مطمئن هستید؟ این عملیات غیرقابل بازگشت است!')) return;

    try {
      await groupsService.deleteGroup(id);
      navigate('/groups', { replace: true });
    } catch (err) {
      setError('❌ خطا در حذف گروه');
      console.error(err);
    }
  };

  if (loading) {
    return <div className="loading">🔄 بارگذاری...</div>;
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

  const userId = parseInt(localStorage.getItem('userId'));
  const isAdmin = group.creator_id === userId;

  return (
    <div className="group-settings-container">
      <div className="settings-header">
        <h1>⚙️ تنظیمات گروه</h1>
        <button className="btn-back" onClick={() => navigate('/groups')}>← برگشت</button>
      </div>

      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}

      {/* اطلاعات گروه */}
      <div className="settings-section">
        <h2>📋 اطلاعات گروه</h2>
        {!editMode ? (
          <div className="group-info">
            <p><strong>📁 نام:</strong> {group.name}</p>
            <p><strong>📝 توضیحات:</strong> {group.description || 'بدون توضیح'}</p>
            <p><strong>👤 مدیر:</strong> {group.creator?.username || 'نامعلوم'}</p>
            {isAdmin && (
              <button className="btn-edit" onClick={() => setEditMode(true)}>
                ✏️ ویرایش اطلاعات
              </button>
            )}
          </div>
        ) : (
          <form onSubmit={handleUpdateGroup} className="edit-form">
            <div className="form-group">
              <label>📁 نام گروه</label>
              <input
                type="text"
                value={editName}
                onChange={(e) => setEditName(e.target.value)}
                required
              />
            </div>
            <div className="form-group">
              <label>📝 توضیحات</label>
              <textarea
                value={editDescription}
                onChange={(e) => setEditDescription(e.target.value)}
                rows="3"
              />
            </div>
            <div className="form-actions">
              <button type="submit" className="btn-primary">💾 ذخیره</button>
              <button type="button" className="btn-secondary" onClick={() => setEditMode(false)}>
                ❌ انصراف
              </button>
            </div>
          </form>
        )}
      </div>

      {/* مدیریت اعضا */}
      <div className="settings-section">
        <h2>👥 مدیریت اعضا ({group.members?.length || 0})</h2>

        {isAdmin && (
          <div className="add-member-section">
            <h3>➕ اضافه کردن عضو جدید</h3>
            <div className="search-box">
              <input
                type="text"
                placeholder="نام کاربری یا ایمیل را جستجو کنید..."
                value={newMemberUsername}
                onChange={(e) => {
                  setNewMemberUsername(e.target.value);
                  handleSearchMembers(e.target.value);
                }}
              />
              {searching && <div className="searching">🔍 در حال جستجو...</div>}
            </div>

            {searchResults.length > 0 && (
              <div className="search-results">
                {searchResults.map(user => (
                  <div key={user.id} className="search-result">
                    <span>👤 {user.username} ({user.email})</span>
                    <button
                      type="button"
                      className="btn-add"
                      onClick={() => handleAddMember(user.id)}
                    >
                      ➕ اضافه
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {group.members && group.members.length > 0 ? (
          <div className="members-grid">
            {group.members.map(member => {
              const isMemberAdmin = member.role === 'admin';
              
              return (
                <div key={member.id} className="member-card">
                  <div className="member-header">
                    <span className="member-name">
                      👤 {member.user?.username || member.username || 'کاربر'}
                    </span>
                    <span className={`member-role ${member.role}`}>
                      {isMemberAdmin ? '👑 مدیر' : member.role === 'member' ? '👥 عضو' : '⏳ معلق'}
                    </span>
                  </div>
                  <p className="member-email">{member.user?.email || member.email || 'بدون ایمیل'}</p>
                  {member.accepted === false && (
                    <p className="pending-text">⏳ در انتظار تایید دعوت</p>
                  )}
                  
                  {isAdmin && !isMemberAdmin && (
                    <div style={{ display: 'flex', gap: '8px', marginTop: '12px' }}>
                      <button
                        className="btn-remove"
                        onClick={() => handleRemoveMember(member.user_id || member.id)}
                        style={{ flex: 1 }}
                      >
                        🗑️ حذف
                      </button>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        ) : (
          <p className="no-members">❌ هیچ عضوی در گروه نیست</p>
        )}
      </div>

      {/* آمار گروه */}
      <div className="settings-section">
        <h2>📊 آمار گروه</h2>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
          <div style={{ padding: '12px', backgroundColor: '#f0f9ff', borderRadius: '8px' }}>
            <p style={{ margin: 0, fontSize: '12px', color: '#6b7280' }}>👥 تعداد اعضا</p>
            <p style={{ margin: '8px 0 0 0', fontSize: '24px', fontWeight: 'bold', color: '#3b82f6' }}>
              {group.members?.length || 0}
            </p>
          </div>
          <div style={{ padding: '12px', backgroundColor: '#f0fdf4', borderRadius: '8px' }}>
            <p style={{ margin: 0, fontSize: '12px', color: '#6b7280' }}>📝 تسک‌ها</p>
            <p style={{ margin: '8px 0 0 0', fontSize: '24px', fontWeight: 'bold', color: '#10b981' }}>
              {group.tasks?.length || 0}
            </p>
          </div>
        </div>
      </div>

      {/* حذف گروه */}
      {isAdmin && (
        <div className="settings-section danger-zone">
          <h2>⚠️ منطقه خطرناک</h2>
          <p style={{ color: '#6b7280', marginBottom: '16px' }}>
            این عملیات غیرقابل بازگشت است.
          </p>
          <button
            className="btn-danger"
            onClick={handleDeleteGroup}
          >
            🗑️ حذف گروه
          </button>
        </div>
      )}
    </div>
  );
}

export default GroupSettings;
