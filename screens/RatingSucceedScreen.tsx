import React from 'react';
import { View, StyleSheet, SafeAreaView } from 'react-native';
import { Button, Text, Card, Paragraph } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import { NativeStackScreenProps } from '@react-navigation/native-stack';

type RootStackParamList = {
  RatingSucceedScreen: undefined;
  NotificationScreen: undefined;
  HomeScreen: undefined;
};

type ScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'RatingSucceedScreen' | 'HomeScreen' 
>;

type NotificationsProps = NativeStackScreenProps<RootStackParamList, 'NotificationScreen'>;

const RatingSucceedScreen: React.FC<NotificationsProps> = ({ navigation }: NotificationsProps) => {
  const use_navigation = useNavigation<ScreenNavigationProp>();

  const GoToHome = () => {
    navigation.navigate('NotificationScreen'); // first let notification stack get back to NotificationScreen
    use_navigation.navigate('HomeScreen'); // then jump to HomeScreen in main stack
  };
  return (
    <SafeAreaView  style={styles.container}> 
      <View style={styles.headerContainer}>
        <Text style={styles.header}>GiveGetGo</Text>
      </View>
      <Card style={styles.card}>
        <Card.Content>
          <Paragraph style={styles.paragraph}>
            Rate submitted!
          </Paragraph>
        </Card.Content>
        <Card.Actions style={styles.cardActions}>
          <Button style={styles.button} mode="contained" onPress={GoToHome}>
            Home
          </Button>
        </Card.Actions>
      </Card>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    marginTop: 50,
    alignItems: 'center',
  },
  header: {
    fontSize: 20,
    fontWeight: 'bold',
    padding: 16,
    alignItems: 'center',
  },
  headerContainer: {
    flexDirection: 'row', // Aligns items in a row
    alignItems: 'center', // Centers items vertically
    paddingLeft: 10, // Adds padding to the left of the avatar
    paddingRight: 10, // Adds padding to the right side
  },
  card: {
    width: '100%',
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  paragraph: {
    textAlign: 'center',
    fontWeight: 'bold',
    fontSize: 16,
    marginBottom: 25,
  },
  button: {
    position: 'absolute', 
    left: 120,
    right: 120, //position, left, right together controls the button's length and horizontal location
    alignSelf: 'center', 
  },
  cardActions: {
    justifyContent: 'center', 
    alignItems: 'center',
    padding: 15,
    width: '100%' // This ensures the actions container fills the width of the card
  },
});

export default RatingSucceedScreen;
