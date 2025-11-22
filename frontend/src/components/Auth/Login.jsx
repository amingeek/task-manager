import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import api from '../../services/api';
import './Auth.css';

export default function Login() {
    const [form, setForm] = useState({ username: '', password: '' });
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        try {
            const response = await api.post('/login', form);
            localStorage.setItem('token', response.data.data.token);
            navigate('/');
        } catch (error) {
            alert('خطا در ورود: ' + (error.response?.data?.error || 'نام کاربری یا رمز عبور اشتباه است'));
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="auth-container">
            <div className="auth-card">
                <div className="auth-header">
                    <h1>ورود به حساب کاربری</h1>
                    <p>خوش آمدید! لطفا وارد شوید</p>
                </div>

                <form onSubmit={handleSubmit} className="auth-form">
                    <div className="form-group">
                        <input
                            type="text"
                            placeholder="نام کاربری"
                            value={form.username}
                            onChange={(e) => setForm({...form, username: e.target.value})}
                            className="form-input"
                            required
                        />
                    </div>

                    <div className="form-group">
                        <input
                            type="password"
                            placeholder="رمز عبور"
                            value={form.password}
                            onChange={(e) => setForm({...form, password: e.target.value})}
                            className="form-input"
                            required
                        />
                    </div>

                    <button
                        type="submit"
                        className="btn btn-primary auth-btn"
                        disabled={loading}
                    >
                        {loading ? 'در حال ورود...' : 'ورود'}
                    </button>
                </form>

                <div className="auth-footer">
                    <p>حساب کاربری ندارید؟</p>
                    <Link to="/register" className="auth-link">
                        ثبت‌نام کنید
                    </Link>
                </div>
            </div>
        </div>
    );
}