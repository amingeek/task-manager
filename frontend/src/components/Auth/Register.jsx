import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import api from '../../services/api';
import './Auth.css';

export default function Register() {
    const [form, setForm] = useState({
        username: '',
        email: '',
        password: '',
        confirmPassword: ''
    });
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (form.password !== form.confirmPassword) {
            alert('رمز عبور و تکرار آن مطابقت ندارند');
            return;
        }

        if (form.password.length < 6) {
            alert('رمز عبور باید حداقل ۶ کاراکتر باشد');
            return;
        }

        setLoading(true);
        try {
            const response = await api.post('/register', {
                username: form.username,
                email: form.email,
                password: form.password
            });

            localStorage.setItem('token', response.data.data.token);
            alert('ثبت‌نام موفقیت‌آمیز بود!');
            navigate('/');
        } catch (error) {
            alert('خطا در ثبت‌نام: ' + (error.response?.data?.error || 'Unknown error'));
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="auth-container">
            <div className="auth-card">
                <div className="auth-header">
                    <h1>ثبت‌نام در تسک منیجر</h1>
                    <p>حساب کاربری جدید ایجاد کنید</p>
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
                            type="email"
                            placeholder="ایمیل"
                            value={form.email}
                            onChange={(e) => setForm({...form, email: e.target.value})}
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
                            minLength="6"
                        />
                    </div>

                    <div className="form-group">
                        <input
                            type="password"
                            placeholder="تکرار رمز عبور"
                            value={form.confirmPassword}
                            onChange={(e) => setForm({...form, confirmPassword: e.target.value})}
                            className="form-input"
                            required
                        />
                    </div>

                    <button
                        type="submit"
                        className="btn btn-primary auth-btn"
                        disabled={loading}
                    >
                        {loading ? 'در حال ثبت‌نام...' : 'ثبت‌نام'}
                    </button>
                </form>

                <div className="auth-footer">
                    <p>حساب کاربری دارید؟</p>
                    <Link to="/login" className="auth-link">
                        ورود به حساب کاربری
                    </Link>
                </div>
            </div>
        </div>
    );
}