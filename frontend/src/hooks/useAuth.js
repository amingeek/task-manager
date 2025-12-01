// frontend/src/hooks/useAuth.js
import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import authService from '../services/auth';

export function useAuth() {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [token, setToken] = useState(null);
    const navigate = useNavigate();

    useEffect(() => {
        const checkAuth = async () => {
            try {
                const storedToken = localStorage.getItem('token');
                if (!storedToken) {
                    setLoading(false);
                    return;
                }

                setToken(storedToken);

                // دریافت اطلاعات کاربر
                const response = await authService.getCurrentUser();
                if (response.data && response.data.data) {
                    setUser(response.data.data);
                } else {
                    setUser(response.data);
                }
            } catch (err) {
                console.error('Auth check failed:', err);
                localStorage.removeItem('token');
                localStorage.removeItem('userId');
                setToken(null);
                setUser(null);
            } finally {
                setLoading(false);
            }
        };

        checkAuth();
    }, []);

    const login = async (username, password) => {
        setLoading(true);
        setError(null);
        try {
            const response = await authService.login(username, password);
            const responseData = response.data || response;

            if (responseData.data && responseData.data.token) {
                const { token: newToken, user: userData } = responseData.data;

                setUser(userData);
                setToken(newToken);
                localStorage.setItem('token', newToken);
                localStorage.setItem('userId', userData?.id || userData?.ID);

                navigate('/dashboard', { replace: true });
                return responseData;
            } else {
                throw new Error('Invalid response format');
            }
        } catch (err) {
            const errorMsg = err.response?.data?.error || err.message || 'خطا در ورود';
            setError(errorMsg);
            throw new Error(errorMsg);
        } finally {
            setLoading(false);
        }
    };

    const register = async (username, email, password) => {
        setLoading(true);
        setError(null);
        try {
            const response = await authService.register(username, email, password);
            const responseData = response.data || response;

            if (responseData.data && responseData.data.token) {
                const { token: newToken, user: userData } = responseData.data;

                setUser(userData);
                setToken(newToken);
                localStorage.setItem('token', newToken);
                localStorage.setItem('userId', userData?.id || userData?.ID);

                navigate('/dashboard', { replace: true });
                return responseData;
            } else {
                throw new Error('Invalid response format');
            }
        } catch (err) {
            const errorMsg = err.response?.data?.error || err.message || 'خطا در ثبت‌نام';
            setError(errorMsg);
            throw new Error(errorMsg);
        } finally {
            setLoading(false);
        }
    };

    const logout = () => {
        authService.logout();
        setUser(null);
        setToken(null);
        localStorage.removeItem('token');
        localStorage.removeItem('userId');
        navigate('/login', { replace: true });
    };

    return {
        user,
        loading,
        error,
        token,
        login,
        register,
        logout,
        isAuthenticated: !!token,
    };
}

export default useAuth;