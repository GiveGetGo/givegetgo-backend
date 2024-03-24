import React, { useState } from 'react';
import { StyleSheet, View, Text, TextInput, TouchableOpacity } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';

// Define the types for your navigation stack
type RootStackParamList = {
    MainScreen: undefined;
};

// Define the type for the navigation prop
type ScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'MainScreen' 
>;

const MainScreen: React.FC = () => {
    return (
        <View style={styles.container}>
        </View>
      );

};

const styles = StyleSheet.create({

});

export default MainScreen;
