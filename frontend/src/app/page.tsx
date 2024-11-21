"use client";

import React, { useEffect, useState } from "react";
import WebApp from "@twa-dev/sdk";
import axios from "axios";
import { WebAppUser } from "@twa-dev/types";
import { logError as LogUtils } from "@/utils/logError";

export default function Home() {
  const [telegramData, setTelegramData] = useState<
    WebAppUser & { added_to_attachment_menu?: boolean; allows_write_to_pm?: boolean } | null
  >(null);
  const [points, setPoints] = useState(0);

  useEffect(() => {
    // Initialize Telegram WebApp
    WebApp.expand();
    WebApp.MainButton.setText("Add Point");
    WebApp.MainButton.show();

    // Get Telegram user data
    const userData = WebApp.initDataUnsafe.user;

    if (userData) {
      setTelegramData(userData);

      // Register or login user
      fetchPoints(userData?.id);
      registerUser(userData.id, userData.username || userData.first_name);

      WebApp.MainButton.onClick(() => {
        addPoints(userData.id);
      });
    }
  }, []);

  const registerUser = async (telegram_id: number, username: string) => {
    try {
      const response = await axios.post(
        `https://miniapp.dandanjan.ir/backend/register`,
        { telegram_id, username }
      );

      setPoints(response.data.points);
    } catch (error) {
      // LogUtils(error, telegram_id, username);
      WebApp.showAlert("Error registering user.");
    }
  };

  const fetchPoints = async (telegram_id: number) => {
    try {
      const response = await axios.get(
        `https://miniapp.dandanjan.ir/backend/points/${telegram_id}?initData=${encodeURIComponent(WebApp.initData)}`
      );
      console.log(`https://miniapp.dandanjan.ir/backend/points/${telegram_id}?initData=${encodeURIComponent(WebApp.initData)}`)
      setPoints(response.data.points);
    } catch (error) {
      // LogUtils(error, telegram_id, telegramData?.username || "unknown");
      WebApp.showAlert("Error fetching points.");
    }
  };

  const addPoints = async (telegram_id: number) => {
    try {
      await axios.post(
        `https://miniapp.dandanjan.ir/backend/points/${telegram_id}?initData=${encodeURIComponent(WebApp.initData)}`
      );      
      fetchPoints(telegram_id); // Refresh points after adding
    } catch (error) {
      // LogUtils(error, telegram_id, telegramData?.username || "unknown");
      WebApp.showAlert("Failed to add points.");
    }
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50">
      <div className="bg-white shadow-lg rounded-lg p-6 max-w-md w-full text-center">
        <h1 className="text-3xl font-bold text-blue-600 mb-4">Mini App By Mehrbod</h1>
        {telegramData ? (
          <>
            <p className="text-lg text-gray-700">
              Welcome, <span className="font-semibold text-blue-600">{telegramData.first_name}</span>!
            </p>
            <p className="mt-2 text-gray-600">
              Your Points: <span className="font-semibold text-green-600">{points}</span>
            </p>
          </>
        ) : (
          <p className="text-gray-500">Loading Telegram Data...</p>
        )}
      </div>
    </div>
  );
}
