import React from 'react';
import { View, StyleSheet, SafeAreaView } from 'react-native';
import { Button, Text, Card, Title, Paragraph } from 'react-native-paper';
import { NativeStackScreenProps } from '@react-navigation/native-stack';

type RootStackParamList = {
  HomeScreen: undefined;
  PostRequestSucceedScreen: undefined; 
};

type HomeScreenProps = NativeStackScreenProps<RootStackParamList, 'HomeScreen'>;

const PostRequestSucceedScreen: React.FC<HomeScreenProps> = ({ navigation }: HomeScreenProps) => {
  const GetContact = () => {
    navigation.navigate('HomeScreen');
  };

  return (
    <SafeAreaView  style={styles.container}> 
      <View style={styles.headerContainer}>
        <Text style={styles.header}>GiveGetGo</Text>
      </View>
      <Card style={styles.card}>
        <Card.Content>
          <Title style={styles.title}>Congratulations!</Title>
          <Paragraph style={styles.paragraph_firstline}>Your request has been submitted. </Paragraph>
          <Paragraph style={styles.paragraph}>         
            You will be able to contact <Paragraph style={styles.boldText}>Jimmy Ho </Paragraph>once a match has been established.
          </Paragraph>
        </Card.Content>
        <Card.Actions style={styles.cardActions}>
          <Button style={styles.button} mode="contained" onPress={GetContact}>
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
  title: {
    textAlign: 'center',
    fontWeight: 'bold',
    marginVertical: 3,
  },
  paragraph: {
    textAlign: 'center',
    fontSize: 16,
    marginBottom: 20,
  },
  paragraph_firstline: {
    textAlign: 'center',
    fontSize: 16,
    marginBottom: 0,
  },
  boldText: {
    fontWeight: 'bold',
    textAlign: 'center',
    fontSize: 16,
    marginBottom: 10,
  },
  button: {
    position: 'absolute', 
    left: 80,
    right: 80, //position, left, right together controls the button's length and horizontal location
    alignSelf: 'center', 
  },
  cardActions: {
    justifyContent: 'center', 
    alignItems: 'center',
    padding: 15,
  },
});

export default PostRequestSucceedScreen;
