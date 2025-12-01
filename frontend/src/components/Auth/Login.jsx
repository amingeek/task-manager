// frontend/src/components/Auth/Login.jsx
import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import './Auth.css';

function Login() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const { login, loading } = useAuth();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        if (!username.trim() || !password.trim()) {
            setError('نام کاربری و رمز عبور ضروری هستند');
            return;
        }

        try {
            await login(username, password);
        } catch (err) {
            setError(err.message || 'خطا در ورود');
        }
    };

    return (
        <div className="auth-container">
            <div className="auth-card">
                <div className="auth-header">
                    <h1>🔐 ورود به سیستم</h1>
                    <p>به تسک منیجر خوش آمدید</p>
                </div>

                {error && (
                    <div className="error-message">
                        {error}
                    </div>
                )}

                <form onSubmit={handleSubmit} className="auth-form">
                    <div className="form-group">
                        <label htmlFor="username">نام کاربری</label>
                        <input
                            id="username"
                            type="text"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            placeholder="نام کاربری خود را وارد کنید"
                            disabled={loading}
                            className="form-input"
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="password">رمز عبور</label>
                        <input
                            id="password"
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            placeholder="رمز عبور خود را وارد کنید"
                            disabled={loading}
                            className="form-input"
                        />
                    </div>

                    <button
                        type="submit"
                        className="btn btn-primary auth-btn"
                        disabled={loading}
                    >
                        {loading ? '⏳ در حال ورود...' : '✅ ورود به سیستم'}
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

export default Login;