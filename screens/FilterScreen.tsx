import React, { useState } from 'react';
import { StyleSheet, View, Text, TextInput, TouchableOpacity } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';

// Define the types for your navigation stack
type RootStackParamList = {
  SignUpScreen: undefined;
  CheckEmailScreen: undefined;
  LoginScreen: undefined;
};

// Define the type for the navigation prop
type ScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'SignUpScreen' | 'CheckEmailScreen' | 'LoginScreen'
>;

const ProfileScreen: React.FC = () => {
  return (
    <View style={styles.container}>
      <Text style={styles.titleText}>Filter</Text>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 16,
    backgroundColor: '#fff',
  },
  titleText: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 48,
  },
  input: {
    height: 40,
    width: '100%',
    borderColor: 'gray',
    borderWidth: 1,
    marginBottom: 16,
    paddingHorizontal: 8,
  },
  button: {
    backgroundColor: 'gray',
    padding: 10,
    borderRadius: 5,
    width: '100%',
    alignItems: 'center',
  },
  buttonText: {
    color: '#fff',
  },
  linkText: {
    color: 'blue',
    marginTop: 16,
  },
});

export default ProfileScreen;
