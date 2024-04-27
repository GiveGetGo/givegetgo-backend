import React, { useState, useEffect } from 'react';
import { View, StyleSheet, Text, SafeAreaView } from 'react-native';
import { Button } from 'react-native-paper';
import { MaterialCommunityIcons } from '@expo/vector-icons';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';

type RootStackParamList = {
  ConfirmationScreen: undefined;
  LoginScreen: undefined;
};

type LoginScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'ConfirmationScreen' | 'LoginScreen'
>;

const ConfirmationScreen: React.FC = () => {
  const navigation = useNavigation<LoginScreenNavigationProp>();
  const [email, setEmail] = useState<string>('email@example.com'); // Placeholder for actual email state

  useEffect(() => {
    // Fetch the email from the backend
    // ...
  }, []);

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.headerContainer}>
          <Text style={styles.header}>GiveGetGo</Text>
      </View>
      <MaterialCommunityIcons name="check-circle" size={100} color="black" />
      <Text style={styles.confirmedText}>Confirmed</Text>
      <Text style={styles.emailText}>{email} has been confirmed</Text>
      <Button mode="contained" onPress={() => navigation.navigate('LoginScreen')} style={styles.button}>
        Home
      </Button>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#fff',
  },
  headerContainer: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    zIndex: 1,
    alignItems: 'center',
    paddingTop: 50, // Adjust for your header's height
  },
  header: {
    fontSize: 22,
    fontWeight: '600',
    color: '#444444',
  },
  confirmedText: {
    fontSize: 24,
    fontWeight: 'bold',
    marginVertical: 16,
  },
  emailText: {
    fontSize: 16,
    marginBottom: 48,
    textAlign: 'center',
  },
  button: {
    paddingHorizontal: 12,
    paddingVertical: 8,
  },
});

export default ConfirmationScreen;
